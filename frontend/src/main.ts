import "./style.css";
import logoUrl from "./assets/logo.webp";

// ---- Wails bridge (window.go / window.runtime are injected at runtime) ----
declare global {
  interface Window {
    go?: any;
    runtime?: any;
  }
}
const App = () => window.go?.main?.App;
const rt = () => window.runtime;

type Snapshot = {
  serverId: number; serverName: string;
  locationId: string; locationName: string;
  playerName: string; inGame: boolean;
  hasEncryptedData: boolean; listening: boolean;
};
type Stats = {
  queued: number; marketOrders: number; marketHistory: number;
  goldPrices: number; emv: number; failed: number;
};
type Feed = { time: string; kind: string; detail: string; count: number };
type Log = { time: string; text: string };
type Config = {
  ingestBaseUrl: string; authBaseUrl: string; requirePow: boolean; packetFilter: string;
  startInTray: boolean; closeToTray: boolean;
  uploadTrades: boolean; uploadMails: boolean; uploadGathering: boolean;
  uploadCombat: boolean; uploadLoot: boolean; uploadParty: boolean;
  uploadSpecs: boolean;
  itemsUrl: string; captureDevice: string;
  marketOrdersTopic: string; marketHistoriesTopic: string; goldPricesTopic: string;
  networkStartDelaySecs: number; idleMinutes: number; idleCheckMinutes: number;
};

let snapshot: Snapshot = {
  serverId: 0, serverName: "", locationId: "", locationName: "",
  playerName: "", inGame: false, hasEncryptedData: false, listening: false,
};
let stats: Stats = { queued: 0, marketOrders: 0, marketHistory: 0, goldPrices: 0, emv: 0, failed: 0 };
type UserStats = { trades: number; mails: number; gathering: number; dungeon: number; loot: number; party: number; specs: number };
let userStats: UserStats = { trades: 0, mails: 0, gathering: 0, dungeon: 0, loot: 0, party: 0, specs: 0 };
let feed: Feed[] = [];
let logs: Log[] = [];
let config: Config | null = null;
let devices: string[] = [];
let authEnabled = false;
let user: { id: string; username: string; avatar: string } = { id: "", username: "", avatar: "" };
let authWaiting = false;
let authPending: { userCode: string; url: string } | null = null;
let authError = "";
let verifyMsg = ""; let verifyMsgCls = "dim"; // persists across the verify tab's live re-renders
// session start time (unix ms, 0 = inactive) per tracker, sourced from the backend
let sessionStarts: Record<string, number> = { gathering: 0, dungeon: 0, loot: 0 };
// feed-event timestamps per session-kind (for the "start a session?" suggestion)
const feedTimes: Record<string, number[]> = { gathering: [], dungeon: [], loot: [] };
const toastShownAt: Record<string, number> = { gathering: 0, dungeon: 0, loot: 0 };
const toastDone: Record<string, boolean> = { gathering: false, dungeon: false, loot: false };
const FEED_KIND_TO_SESSION: Record<string, string> = { gather: "gathering", dungeon: "dungeon", loot: "loot" };
let route = "dashboard";
let updateInfo: any = null;     // latest version-check result being shown
let updateDismissed = "";       // soft update the user clicked "Later" on

// ---- helpers ----
const el = (html: string): HTMLElement => {
  const t = document.createElement("template");
  t.innerHTML = html.trim();
  return t.content.firstElementChild as HTMLElement;
};
const fmt = (n: number) => n.toLocaleString();
const hhmmss = (iso: string) => {
  const d = new Date(iso);
  if (isNaN(+d)) return "";
  return d.toLocaleTimeString();
};

// ---- render ----
function render() {
  const root = document.getElementById("app")!;
  if (!root.querySelector(".shell")) {
    root.innerHTML = "";
    root.appendChild(buildShell());
  }
  renderSidebarStatus();
  renderContent();
}

function buildShell(): HTMLElement {
  const shell = el(`<div class="shell">
    <aside class="sidebar">
      <div class="brand">
        <img src="${logoUrl}" alt="logo"/>
        <div>
          <div class="name">Albion Market</div>
          <div class="sub">Data Client</div>
        </div>
      </div>
      <nav class="nav" id="nav"></nav>
      <div class="authcard" id="authcard"></div>
      <div class="statuscard" id="statuscard"></div>
    </aside>
    <main class="main">
      <div class="topbar">
        <h2 id="title">Dashboard</h2>
        <div class="spacer"></div>
        <div id="topactions" class="row"></div>
        <button class="btn ghost" id="trayBtn" title="Hide to system tray (keeps running in the background)">⤓ Tray</button>
      </div>
      <div class="content" id="content"></div>
    </main>
  </div>`);

  const nav = shell.querySelector("#nav")!;
  const items: [string, string][] = [
    ["dashboard", "Dashboard"],
    ["feed", "Live Feed"],
    ["trades", "Trades"],
    ["mails", "Mails"],
    ["dungeon", "Dungeon"],
    ["gathering", "Gathering"],
    ["loot", "Loot"],
    ["verify", "Verification"],
    ["logs", "Logs"],
    ["settings", "Settings"],
  ];
  for (const [id, label] of items) {
    const b = el(`<button data-route="${id}">${label}</button>`);
    b.addEventListener("click", () => { route = id; render(); });
    nav.appendChild(b);
  }
  shell.querySelector("#trayBtn")!.addEventListener("click", () => { App()?.HideToTray(); });
  return shell;
}

// avatarSrc resolves a usable avatar image URL. The backend may send `avatar`
// as a full URL, or as a Discord avatar hash (in which case we need the Discord
// snowflake id). If neither yields a URL, callers fall back to an initial.
function avatarSrc(u: { id: string; avatar: string }): string {
  const a = (u.avatar || "").trim();
  if (!a) return "";
  if (/^https?:\/\//.test(a)) return a;
  if (/^\d{16,20}$/.test(String(u.id))) {
    const ext = a.startsWith("a_") ? "gif" : "png";
    return `https://cdn.discordapp.com/avatars/${u.id}/${a}.${ext}?size=64`;
  }
  return "";
}

// mountAvatar fills a wrapper with the avatar image, falling back to a coloured
// initial bubble on missing/failed image (never a broken-image icon).
function mountAvatar(wrap: HTMLElement, u: { id: string; username: string; avatar: string }) {
  const initial = (u.username || "?").charAt(0).toUpperCase();
  const gen = () => { wrap.innerHTML = `<span class="avatar gen">${escapeHtml(initial)}</span>`; };
  const src = avatarSrc(u);
  if (!src) { gen(); return; }
  const img = document.createElement("img");
  img.className = "avatar";
  img.referrerPolicy = "no-referrer";
  img.onerror = gen;
  img.src = src;
  wrap.innerHTML = "";
  wrap.appendChild(img);
}

async function startLogin() {
  authError = ""; authPending = null; authWaiting = true;
  renderAuthCard();
  if (route === "settings") renderContent();
  await App()?.Login();
}

function renderAuthCard() {
  const ac = document.getElementById("authcard");
  if (!ac) return;
  if (!authEnabled) { ac.innerHTML = ""; ac.style.display = "none"; return; }
  ac.style.display = "";
  if (user.id) {
    ac.innerHTML = `<div class="authrow"><span class="avwrap"></span><div class="aname">${escapeHtml(user.username || "Signed in")}</div></div>
      <button class="btn ghost sm" id="logoutBtn">Sign out</button>`;
    mountAvatar(ac.querySelector(".avwrap") as HTMLElement, user);
    ac.querySelector("#logoutBtn")!.addEventListener("click", async () => { await App()?.Logout(); user = { id: "", username: "", avatar: "" }; renderAuthCard(); });
  } else if (authWaiting) {
    const fallback = authPending
      ? `<div class="dim" style="font-size:11px;margin-top:6px">If the browser didn't open, go to<br/><a href="${escapeHtml(authPending.url)}" target="_blank" class="acc">${escapeHtml(authPending.url)}</a><br/>code <b>${escapeHtml(authPending.userCode)}</b></div>`
      : "";
    ac.innerHTML = `<div class="live"><span class="dot warn"></span> Waiting for browser sign-in…</div>${fallback}
      <button class="btn ghost sm" id="cancelLoginBtn" style="margin-top:8px">Cancel</button>`;
    ac.querySelector("#cancelLoginBtn")!.addEventListener("click", () => { authWaiting = false; authPending = null; renderAuthCard(); });
  } else {
    ac.innerHTML = `<button class="btn discord" id="loginBtn">Login with Discord</button>
      <div class="dim" style="font-size:11px;margin-top:6px">Attribute your uploads (optional)</div>
      ${authError ? `<div class="bad" style="font-size:11px;margin-top:6px">${escapeHtml(authError)}</div>` : ""}`;
    ac.querySelector("#loginBtn")!.addEventListener("click", startLogin);
  }
}

function renderSidebarStatus() {
  document.querySelectorAll<HTMLElement>(".nav button").forEach((b) => {
    b.classList.toggle("active", b.dataset.route === route);
  });
  renderAuthCard();
  const sc = document.getElementById("statuscard");
  if (!sc) return;
  const dot = (cls: string) => `<span class="dot ${cls}"></span>`;
  const ok = (b: boolean) => (b ? "ok" : "bad");
  sc.innerHTML = `
    <h4>Status</h4>
    <div class="statusrow">${dot(snapshot.listening ? "ok" : "warn")} Capture ${snapshot.listening ? "running" : "stopped"}</div>
    <div class="statusrow">${dot(ok(snapshot.inGame))} ${snapshot.inGame ? "In game" : "No game traffic"}</div>
    <div class="statusrow">${dot(snapshot.serverId ? "ok" : "warn")} ${snapshot.serverName || "Server unknown"}</div>
    <div class="statusrow">${dot(snapshot.locationId ? "ok" : "warn")} ${snapshot.locationName || snapshot.locationId || "Location unknown"}</div>
    <div class="statusrow">${dot(snapshot.hasEncryptedData ? "bad" : "ok")} ${snapshot.hasEncryptedData ? "Encrypted (no orders)" : "Data readable"}</div>`;
}

function renderContent() {
  const title = document.getElementById("title")!;
  const content = document.getElementById("content")!;
  const actions = document.getElementById("topactions")!;
  actions.innerHTML = "";
  title.textContent = ({ dashboard: "Dashboard", feed: "Live Feed", trades: "Trades", mails: "Mails", dungeon: "Dungeon", gathering: "Gathering", loot: "Loot", verify: "Verification", logs: "Logs", settings: "Settings" } as any)[route];

  if (route === "dashboard") return renderDashboard(content, actions);
  if (route === "feed") return renderFeed(content, actions);
  if (route === "trades") return renderTrades(content, actions);
  if (route === "mails") return renderMails(content, actions);
  if (route === "dungeon") return renderDungeon(content, actions);
  if (route === "gathering") return renderGathering(content, actions);
  if (route === "loot") return renderLoot(content, actions);
  if (route === "verify") return renderVerify(content, actions);
  if (route === "logs") return renderLogs(content, actions);
  if (route === "settings") return renderSettings(content, actions);
}

type Character = { name: string; fame: number; serverId: number; serverName: string };

async function renderVerify(content: HTMLElement, _actions: HTMLElement) {
  if (!authEnabled || !user.id) {
    content.innerHTML = `<div class="empty">Sign in with Discord (sidebar) to verify your character to your Albion Market account.</div>`;
    return;
  }
  const chars: Character[] = (await App()?.GetDetectedCharacters()) || [];
  const charsHTML = chars.length
    ? chars.map((c) =>
        `<tr><td class="mono">${escapeHtml(c.name)}</td><td class="dim">${c.serverName || "—"}</td><td class="mono">${silver(c.fame)} fame</td>
         <td style="text-align:right"><button class="btn sm vbtn" data-name="${escapeHtml(c.name)}" data-srv="${c.serverId}" data-fame="${c.fame}" style="width:auto;padding:0 12px">Verify</button></td></tr>`
      ).join("")
    : `<tr><td colspan="4" class="dim">No character detected yet — log into Albion (or travel between zones) and your character appears here.</td></tr>`;

  content.innerHTML = `
    <div class="panel" style="max-width:680px">
      <h3>Verify a character</h3>
      <div class="panel-body">
        <div class="dim" style="font-size:12px;margin-bottom:12px">
          We detect your in-game character straight from the game's login packet — proving you actually
          play it. Click <b>Verify</b> to send it to your Albion Market account; the result shows below.
        </div>
        <table><thead><tr><th>Detected character</th><th>Server</th><th>Fame</th><th></th></tr></thead>
          <tbody>${charsHTML}</tbody></table>
        <div id="vmsg" class="${verifyMsgCls}" style="font-size:13px;margin-top:14px">${escapeHtml(verifyMsg)}</div>
      </div>
    </div>`;

  content.querySelectorAll<HTMLButtonElement>(".vbtn").forEach((b) => {
    b.addEventListener("click", async () => {
      b.disabled = true; b.textContent = "Verifying…";
      verifyMsg = "Verifying…"; verifyMsgCls = "dim";
      const res = await App()?.SubmitVerification(parseInt(b.dataset.srv!), b.dataset.name!, parseInt(b.dataset.fame!));
      verifyMsg = res?.message || "Done.";
      verifyMsgCls = (res?.status === "verified" || res?.status === "already_verified") ? "good" : ((res?.status && res.status !== "error") ? "bad" : "dim");
      const msg = content.querySelector("#vmsg");
      if (msg) { msg.textContent = verifyMsg; msg.className = verifyMsgCls; }
      b.disabled = false; b.textContent = "Verify";
    });
  });
}

const dur = (s: number) => {
  const h = Math.floor(s / 3600), m = Math.floor((s % 3600) / 60), sec = s % 60;
  return (h ? h + "h " : "") + (m || h ? m + "m " : "") + sec + "s";
};

// Capture + sync tool by design: analytics (totals, value, DPS, silver/hr, trends)
// live on the website, fed by sync. These tabs show only a stopwatch + capture log.

const websiteNote = `<div class="dim" style="margin-top:12px;font-size:12px">Detailed stats, totals and value are shown on the <b>Albion Market</b> website — this client captures the data and syncs it to your account.</div>`;

// renderKindLog shows the raw capture log filtered to the given feed kind(s).
function renderKindLog(host: HTMLElement, kinds: string[], emptyMsg: string) {
  const rows = feed.filter((f) => kinds.includes(f.kind)).slice(-300).reverse();
  if (!rows.length) { host.innerHTML = `<div class="empty">${emptyMsg}</div>`; return; }
  host.innerHTML = `<table><thead><tr><th>Time</th><th>Captured</th><th>Qty</th></tr></thead><tbody>${
    rows.map((f) => `<tr><td class="mono dim">${hhmmss(f.time)}</td><td>${escapeHtml(f.detail)}</td><td class="mono">${fmt(f.count)}</td></tr>`).join("")
  }</tbody></table>`;
}

const hhmmssDur = (ms: number) => {
  const s = Math.max(0, Math.floor(ms / 1000));
  const p = (n: number) => String(n).padStart(2, "0");
  return `${p(Math.floor(s / 3600))}:${p(Math.floor((s % 3600) / 60))}:${p(s % 60)}`;
};

async function toggleSession(key: string, start: boolean) {
  if (start) { await App()?.StartSession(key); sessionStarts[key] = Date.now(); toastDone[key] = true; toastShownAt[key] = 0; }
  else { await App()?.StopSession(key); sessionStarts[key] = 0; }
  renderToasts();
  renderContent();
}

function sessionBar(actions: HTMLElement, key: string) {
  const running = sessionStarts[key] > 0;
  const btn = el(`<button class="btn ${running ? "ghost" : ""}">${running ? "Stop Session" : "Start Session"}</button>`);
  btn.addEventListener("click", () => toggleSession(key, !running));
  actions.appendChild(btn);
}

// renderSessionTab is the shared capture-only view: an independent session
// stopwatch + the raw capture log for one tracker kind. Sessions run concurrently.
function renderSessionTab(content: HTMLElement, actions: HTMLElement, key: string, title: string, logKind: string, emptyMsg: string) {
  sessionBar(actions, key);
  const start = sessionStarts[key];
  const live = start > 0;
  content.innerHTML = `
    <div class="grid cards" style="margin-bottom:16px">
      <div class="card"><div class="label">${title} session ${live ? '<span class="live"><span class="dot ok"></span>running</span>' : ""}</div><div class="value accent">${hhmmssDur(live ? Date.now() - start : 0)}</div></div>
    </div>
    <div class="panel"><h3>${title} capture log</h3><div id="klog"></div></div>
    ${websiteNote}`;
  renderKindLog(content.querySelector("#klog")!, [logKind], emptyMsg);
}

function renderGathering(content: HTMLElement, actions: HTMLElement) {
  renderSessionTab(content, actions, "gathering", "Gathering", "gather", "No gathering captured yet. Harvest or fish in-game and it appears here.");
}

function renderDungeon(content: HTMLElement, actions: HTMLElement) {
  renderSessionTab(content, actions, "dungeon", "Dungeon", "dungeon", "No dungeon activity captured yet. Fame and silver from dungeon runs appear here.");
}

function renderLoot(content: HTMLElement, actions: HTMLElement) {
  renderSessionTab(content, actions, "loot", "Loot", "loot", "No loot captured yet. Items picked up by you and nearby players appear here.");
}

type Mail = { itemId: string; auctionType: number; partialAmount: number; totalAmount: number; unitSilver: number; totalSilver: number; received: string; playerName: string; locationId: number; isSet: boolean };
type Trade = { itemId: string; operation: number; type: number; amount: number; unitSilver: number; qualityLevel: number; playerName: string; dateTime: string; locationId: number };

const silver = (n: number) => Math.round(n).toLocaleString();

async function renderTrades(content: HTMLElement, _actions: HTMLElement) {
  content.innerHTML = `<div class="panel"><h3>Trades</h3><div id="t"></div></div>`;
  const host = content.querySelector("#t")!;
  const rows: Trade[] = (await App()?.GetTrades()) || [];
  if (!rows.length) { host.innerHTML = `<div class="empty">No trades captured yet. Instant buys/sells and filled market orders (from mail) appear here.</div>`; return; }
  host.innerHTML = `<table><thead><tr><th>Time</th><th>Op</th><th>Type</th><th>Item</th><th>Qty</th><th>Unit (silver)</th></tr></thead><tbody>${
    rows.map((t) => `<tr>
      <td class="mono dim">${hhmmss(t.dateTime)}</td>
      <td><span class="tag ${t.operation === 0 ? "buy" : "sell"}">${t.operation === 0 ? "BUY" : "SELL"}</span></td>
      <td class="dim">${t.type === 0 ? "Instant" : "Order"}</td>
      <td class="mono">${escapeHtml(t.itemId)}</td>
      <td class="mono">${fmt(t.amount)}</td>
      <td class="mono">${silver(t.unitSilver)}</td></tr>`).join("")
  }</tbody></table>`;
}

async function renderMails(content: HTMLElement, _actions: HTMLElement) {
  content.innerHTML = `<div class="panel"><h3>Marketplace Mails</h3><div id="m"></div></div>`;
  const host = content.querySelector("#m")!;
  const rows: Mail[] = (await App()?.GetMails()) || [];
  if (!rows.length) { host.innerHTML = `<div class="empty">No mails captured yet. Open your in-game marketplace mailbox to collect sold/bought summaries.</div>`; return; }
  host.innerHTML = `<table><thead><tr><th>Received</th><th>Side</th><th>Item</th><th>Filled</th><th>Unit (silver)</th><th>Total (silver)</th></tr></thead><tbody>${
    rows.map((m) => `<tr>
      <td class="mono dim">${hhmmss(m.received)}</td>
      <td><span class="tag ${m.auctionType === 1 ? "sell" : "buy"}">${m.auctionType === 1 ? "Sold" : "Bought"}</span></td>
      <td class="mono">${escapeHtml(m.itemId || "—")}</td>
      <td class="mono">${fmt(m.partialAmount)}/${fmt(m.totalAmount)}</td>
      <td class="mono">${silver(m.unitSilver)}</td>
      <td class="mono">${silver(m.totalSilver)}</td></tr>`).join("")
  }</tbody></table>`;
}

function captureButtons(actions: HTMLElement) {
  const live = el(`<span class="live">${snapshot.listening ? '<span class="dot ok"></span> live' : ""}</span>`);
  actions.appendChild(live);
  const toggle = el(`<button class="btn ${snapshot.listening ? "ghost" : ""}">${snapshot.listening ? "Stop Capture" : "Start Capture"}</button>`);
  toggle.addEventListener("click", () => { App()?.ToggleCapture(!snapshot.listening); });
  actions.appendChild(toggle);
}

function renderDashboard(content: HTMLElement, actions: HTMLElement) {
  captureButtons(actions);
  const cardsDiv = (label: string, value: string, cls = "") =>
    `<div class="card"><div class="label">${label}</div><div class="value ${cls}">${value}</div></div>`;
  content.innerHTML = `
    <div class="grid cards" style="margin-bottom:16px">
      ${cardsDiv("Market Orders", fmt(stats.marketOrders), "accent")}
      ${cardsDiv("Item History", fmt(stats.marketHistory))}
      ${cardsDiv("Gold Prices", fmt(stats.goldPrices), "good")}
      ${cardsDiv("Est. Values", fmt(stats.emv))}
      ${cardsDiv("Queued", fmt(stats.queued))}
      ${cardsDiv("Failed", fmt(stats.failed), stats.failed ? "bad" : "")}
    </div>
    <div class="section-title">Your data synced ${user.id ? "" : "(sign in to sync)"}</div>
    <div class="grid cards" style="margin-bottom:16px">
      ${cardsDiv("Trades", fmt(userStats.trades))}
      ${cardsDiv("Mails", fmt(userStats.mails))}
      ${cardsDiv("Gathering", fmt(userStats.gathering))}
      ${cardsDiv("Dungeon", fmt(userStats.dungeon))}
      ${cardsDiv("Loot", fmt(userStats.loot))}
      ${cardsDiv("Party", fmt(userStats.party))}
      ${cardsDiv("Specs", fmt(userStats.specs))}
    </div>
    <div class="panel">
      <h3>Recent Captures</h3>
      <div id="recent"></div>
    </div>`;
  renderFeedTable(content.querySelector("#recent")!, feed.slice(-12).reverse());
}

function renderFeed(content: HTMLElement, actions: HTMLElement) {
  captureButtons(actions);
  content.innerHTML = `<div class="panel"><h3>Live Capture Feed</h3><div id="feedtable"></div></div>`;
  renderFeedTable(content.querySelector("#feedtable")!, [...feed].reverse());
}

function renderFeedTable(host: HTMLElement, rows: Feed[]) {
  if (!rows.length) {
    host.innerHTML = `<div class="empty">No captures yet. Start Albion Online (with capture running) to see market data flow in.</div>`;
    return;
  }
  const body = rows.map((f) => `
    <tr>
      <td class="mono dim">${hhmmss(f.time)}</td>
      <td><span class="tag ${f.kind}">${f.kind}</span></td>
      <td>${f.detail}</td>
      <td class="mono">${fmt(f.count)}</td>
    </tr>`).join("");
  host.innerHTML = `<table>
    <thead><tr><th>Time</th><th>Kind</th><th>Detail</th><th>Count</th></tr></thead>
    <tbody>${body}</tbody></table>`;
}

function renderLogs(content: HTMLElement, _actions: HTMLElement) {
  content.innerHTML = `<div class="panel"><h3>Logs</h3><div class="logs" id="loglist"></div></div>`;
  const list = content.querySelector("#loglist")!;
  if (!logs.length) { list.innerHTML = `<div class="dim">No log output yet.</div>`; return; }
  list.innerHTML = logs.slice(-300).map((l) =>
    `<div class="logline"><span class="t">${hhmmss(l.time)}</span>  ${escapeHtml(l.text)}</div>`).join("");
  list.scrollTop = list.scrollHeight;
}

function renderSettings(content: HTMLElement, _actions: HTMLElement) {
  if (!config) { content.innerHTML = `<div class="empty">Loading settings…</div>`; return; }
  const c = config;
  let account: string;
  if (user.id) {
    account = `<div class="row"><span>Signed in as <b>${escapeHtml(user.username || "—")}</b></span>
        <button class="btn ghost sm" id="acctLogout" style="width:auto">Sign out</button></div>`;
  } else if (authWaiting) {
    account = `<div class="live"><span class="dot warn"></span> Waiting for browser sign-in…</div>`;
  } else {
    account = `<button class="btn discord" id="acctLogin" style="width:auto;padding:0 16px">Login with Discord</button>`;
  }

  const dataTypes: [string, string, string, boolean][] = [
    ["uploadTrades", "Trades", "Your instant buys/sells and filled market orders.", c.uploadTrades],
    ["uploadMails", "Marketplace mails", "Sold/bought results from your marketplace mailbox.", c.uploadMails],
    ["uploadGathering", "Gathering & fishing", "Resources & fish you gather, with session value.", c.uploadGathering],
    ["uploadCombat", "Dungeon", "Fame, silver and loot summaries from your dungeon runs.", c.uploadCombat],
    ["uploadLoot", "Loot", "Items looted by you and nearby players.", c.uploadLoot],
    ["uploadParty", "Party", "Your current party members.", c.uploadParty],
    ["uploadSpecs", "Character specs", "Your crafting/gathering/combat mastery & specialization levels (sent on login).", c.uploadSpecs],
  ];
  const dataToggles = dataTypes.map(([id, name, desc, on]) => `
    <label class="datatoggle" for="${id}">
      <input type="checkbox" id="${id}" ${on ? "checked" : ""}/>
      <div><div class="dt-name">${name}</div><div class="dt-desc">${desc}</div></div>
    </label>`).join("");

  content.innerHTML = `
    <div class="grid settings-grid">
      <div class="panel settings-wide">
        <h3>Your data sync</h3>
        <div class="panel-body">
          <div class="row" style="justify-content:space-between;margin-bottom:14px">
            <div>${account}</div>
          </div>
          <div class="dim" style="font-size:12px;margin-bottom:12px">
            Choose which of <b>your own</b> captured data syncs to your Albion Market account.
            Only uploads while signed in — it powers your history, portfolio &amp; stats on the website.
            Market prices everyone benefits from are always shared separately.
          </div>
          <div class="datatoggles">${dataToggles}</div>
        </div>
      </div>

      <div class="panel">
        <h3>Network adapter</h3>
        <div class="panel-body">
          <div class="field">
            <label>Leave on <b>All</b> unless capture isn't working — then pick the adapter you play on.</label>
            <select id="captureDevice"></select>
          </div>
        </div>
      </div>

      <div class="panel">
        <h3>Window &amp; tray</h3>
        <div class="panel-body">
          <div class="field switch">
            <input type="checkbox" id="startInTray" ${c.startInTray ? "checked" : ""}/>
            <label for="startInTray">Start minimized in the system tray</label>
          </div>
          <div class="field switch">
            <input type="checkbox" id="closeToTray" ${c.closeToTray ? "checked" : ""}/>
            <label for="closeToTray">Closing the window hides to tray (keeps running in the background)</label>
          </div>
        </div>
      </div>
    </div>
    <div class="row" style="margin-top:16px">
      <button class="btn" id="saveBtn">Save Settings</button>
      <span class="dim" id="saveMsg"></span>
    </div>`;

  const sel = content.querySelector("#captureDevice") as HTMLSelectElement;
  const opts = ["", ...devices];
  sel.innerHTML = opts.map((d) => `<option value="${escapeHtml(d)}" ${d === c.captureDevice ? "selected" : ""}>${d ? escapeHtml(d) : "All adapters"}</option>`).join("");

  content.querySelector("#acctLogin")?.addEventListener("click", startLogin);
  content.querySelector("#acctLogout")?.addEventListener("click", async () => {
    await App()?.Logout(); user = { id: "", username: "", avatar: "" }; renderContent(); renderAuthCard();
  });

  content.querySelector<HTMLButtonElement>("#saveBtn")!.addEventListener("click", async () => {
    const chk = (id: string) => (content.querySelector("#" + id) as HTMLInputElement).checked;
    const next: Config = {
      ...c,
      startInTray: chk("startInTray"),
      closeToTray: chk("closeToTray"),
      captureDevice: (content.querySelector("#captureDevice") as HTMLSelectElement).value,
      uploadTrades: chk("uploadTrades"),
      uploadMails: chk("uploadMails"),
      uploadGathering: chk("uploadGathering"),
      uploadCombat: chk("uploadCombat"),
      uploadLoot: chk("uploadLoot"),
      uploadParty: chk("uploadParty"),
      uploadSpecs: chk("uploadSpecs"),
    };
    const err = await App()?.SaveConfig(next);
    const msg = content.querySelector("#saveMsg")!;
    msg.textContent = err ? "Error: " + err : "Saved.";
    config = next;
  });
}

const TOAST_TTL = 5 * 60 * 1000;          // a suggestion stays ~5 min
const TOAST_RATE = 5;                      // events in the last 60s to suggest

// renderToasts draws one "start a session?" suggestion per active kind.
function renderToasts() {
  let host = document.getElementById("toasts");
  if (!host) { host = document.createElement("div"); host.id = "toasts"; document.body.appendChild(host); }
  const active = ["gathering", "dungeon", "loot"].filter((k) => toastShownAt[k] > 0);
  host.innerHTML = active.map((k) => `
    <div class="toast">
      <div>
        <div class="toast-title">Lots of ${k} activity</div>
        <div class="toast-desc">Start a ${k} session to sync it to your account?</div>
      </div>
      <div class="row" style="gap:6px;margin-top:8px">
        <button class="btn sm" data-act="start" data-k="${k}" style="width:auto;padding:0 12px">Start session</button>
        <button class="btn ghost sm" data-act="dismiss" data-k="${k}" style="width:auto;padding:0 12px">Not now</button>
      </div>
    </div>`).join("");
  host.querySelectorAll<HTMLButtonElement>("button[data-act]").forEach((b) => {
    const k = b.dataset.k!;
    b.addEventListener("click", () => {
      if (b.dataset.act === "start") { toggleSession(k, true); }
      else { toastShownAt[k] = 0; toastDone[k] = true; renderToasts(); }
    });
  });
}

// tick refreshes session state and decides whether to surface a suggestion.
async function tick() {
  const s = await App()?.GetSessions();
  if (s) sessionStarts = { gathering: s.gathering || 0, dungeon: s.dungeon || 0, loot: s.loot || 0 };
  const now = Date.now();
  let changed = false;
  for (const k of ["gathering", "dungeon", "loot"]) {
    feedTimes[k] = feedTimes[k].filter((t) => now - t < 60000);
    if (toastShownAt[k] > 0 && now - toastShownAt[k] > TOAST_TTL) { toastShownAt[k] = 0; toastDone[k] = true; changed = true; }
    if (sessionStarts[k] > 0) { toastShownAt[k] = 0; } // active -> no suggestion
    if (feedTimes[k].length >= TOAST_RATE && sessionStarts[k] === 0 && toastShownAt[k] === 0 && !toastDone[k]) {
      toastShownAt[k] = now; changed = true;
    }
  }
  if (changed) renderToasts();
  if (["dungeon", "gathering", "loot", "verify"].includes(route)) renderContent();
}

function escapeHtml(s: string) {
  return s.replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c]!));
}

// ---- wire events + initial load ----
// showUpdate decides whether the version-check result warrants the modal.
function showUpdate(r: any) {
  if (!r) return;
  const required = !!r.updateRequired;
  if (!(required || (r.available && !r.upToDate))) return;
  if (!required && updateDismissed === r.latest) return; // soft update already dismissed
  updateInfo = r;
  renderUpdateModal();
}

function renderUpdateModal() {
  let host = document.getElementById("updateModal");
  const r = updateInfo;
  if (!r) { host?.remove(); return; }
  const required = !!r.updateRequired;
  if (!host) { host = document.createElement("div"); host.id = "updateModal"; document.body.appendChild(host); }
  host.className = "modal-overlay";
  host.innerHTML = `
    <div class="modal">
      <h3>${required ? "Update required" : "Update available"}</h3>
      <p>${required
        ? `Your version (${escapeHtml(r.current || "")}) is no longer supported. Please update to keep using the client.`
        : `A newer version (${escapeHtml(r.latest || "")}) is available — you're on ${escapeHtml(r.current || "")}.`}</p>
      ${r.notes ? `<p class="dim">${escapeHtml(r.notes)}</p>` : ""}
      <div class="row" style="gap:8px;margin-top:16px;justify-content:flex-end">
        ${required ? "" : `<button class="btn ghost sm" id="updLater" style="width:auto;padding:0 14px">Later</button>`}
        <button class="btn sm" id="updGet" style="width:auto;padding:0 14px">Download</button>
      </div>
    </div>`;
  host.querySelector("#updGet")!.addEventListener("click", () => { App()?.OpenUpdateDownload(); });
  host.querySelector("#updLater")?.addEventListener("click", () => {
    updateDismissed = r.latest || ""; updateInfo = null; renderUpdateModal();
  });
}

async function init() {
  render();
  const r = rt();
  if (r?.EventsOn) {
    r.EventsOn("state", (s: Snapshot) => { snapshot = s; renderSidebarStatus(); if (route === "dashboard") renderContent(); });
    r.EventsOn("stats", (s: Stats) => { stats = s; if (route === "dashboard") renderContent(); });
    r.EventsOn("userstats", (s: UserStats) => { userStats = s; if (route === "dashboard") renderContent(); });
    r.EventsOn("feed", (f: Feed) => {
      feed.push(f); if (feed.length > 500) feed.shift();
      const sk = FEED_KIND_TO_SESSION[f.kind]; if (sk) feedTimes[sk].push(Date.now());
      if (route === "feed" || route === "dashboard") renderContent();
    });
    r.EventsOn("log", (l: Log) => { logs.push(l); if (logs.length > 1000) logs.shift(); if (route === "logs") renderContent(); });
    r.EventsOn("auth", (u: any) => { user = u || { id: "", username: "", avatar: "" }; authWaiting = false; authPending = null; authError = ""; renderAuthCard(); if (route === "settings") renderContent(); });
    r.EventsOn("authPending", (p: any) => { authPending = { userCode: p.userCode, url: p.url }; if (authWaiting) renderAuthCard(); });
    r.EventsOn("authError", (msg: string) => { authWaiting = false; authPending = null; authError = msg || "Login failed"; renderAuthCard(); });
    r.EventsOn("update", (u: any) => showUpdate(u));
  }
  // 1s tick: refresh session state, drive the session-suggestion toasts, and
  // keep the tracker tab's timer/log live.
  setInterval(tick, 1000);

  const a = App();
  if (a) {
    try {
      snapshot = await a.GetSnapshot();
      stats = await a.GetStats();
      feed = (await a.GetFeed()) || [];
      logs = (await a.GetLogs()) || [];
      config = await a.GetConfig();
      devices = (await a.GetDevices()) || [];
      authEnabled = (await a.AuthEnabled()) || false;
      user = (await a.GetUser()) || user;
      userStats = (await a.GetUserStats()) || userStats;
      const ss = await a.GetSessions(); if (ss) sessionStarts = { gathering: ss.gathering || 0, dungeon: ss.dungeon || 0, loot: ss.loot || 0 };
    } catch (e) { /* running outside wails */ }
  }
  render();
}

init();
