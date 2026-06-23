package main

import "net/http"

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

const indexHTML = `<!doctype html><html><head><meta charset="utf-8"><title>Albion Packet Inspector</title>
<style>
 :root{--bg0:#161616;--bg1:#1b1b1b;--bg2:#25282b;--bg3:#2d3034;--bd:#383c41;--bds:#2c2f33;--tx:#fff;--txm:#a0a0a0;--txd:#6f757c;--acc:#5865f2;--good:#43b581;--bad:#f04747;--warn:#faa61a}
 *{box-sizing:border-box} html,body{margin:0;height:100%}
 body{background:var(--bg0);color:var(--tx);font:13px Roboto,system-ui,sans-serif;display:flex;flex-direction:column;height:100vh}
 header{background:var(--bg1);border-bottom:1px solid var(--bds);padding:9px 14px;display:flex;gap:12px;align-items:center;flex-wrap:wrap}
 header h1{font-size:14px;margin:0;font-weight:600}
 .stat{color:var(--txd);font-size:12px}.stat b{color:var(--tx)} .unkc{color:var(--bad)}
 .spacer{flex:1}
 input[type=text]{background:var(--bg3);border:1px solid var(--bd);border-radius:6px;color:var(--tx);padding:6px 10px;font:inherit;width:200px}
 label.chk{display:inline-flex;gap:5px;align-items:center;color:var(--txm);font-size:12px;cursor:pointer;user-select:none}
 button{background:var(--bg2);border:1px solid var(--bd);color:var(--txm);border-radius:6px;padding:6px 11px;font:inherit;cursor:pointer}
 button.on{background:var(--acc);color:#fff;border-color:var(--acc)}
 button.mark{background:var(--warn);color:#1b1b1b;border-color:var(--warn);font-weight:600}
 .wrap{flex:1;overflow:auto;position:relative}
 table{width:100%;border-collapse:collapse}
 th,td{text-align:left;padding:5px 12px;border-bottom:1px solid var(--bds);white-space:nowrap}
 th{position:sticky;top:0;background:var(--bg1);color:var(--txd);font-size:11px;text-transform:uppercase;z-index:1}
 tbody tr.r{cursor:pointer} tbody tr.r:hover{background:var(--bg3)}
 tr.mark td{background:rgba(250,166,26,.14);color:var(--warn);font-weight:600;text-align:center;border-top:2px solid var(--warn)}
 .mono{font-family:ui-monospace,Consolas,monospace}
 .badge{padding:1px 7px;border-radius:999px;font-size:11px;font-weight:600}
 .event{background:rgba(88,101,242,.18);color:#9fb0ff}
 .request{background:rgba(250,166,26,.16);color:var(--warn)}
 .response{background:rgba(67,181,129,.16);color:var(--good)}
 .unk{color:var(--bad)}
 .dim{color:var(--txd)}
 pre{margin:0;padding:10px 16px 14px 40px;background:var(--bg2);color:#cdd3da;font:12px ui-monospace,Consolas,monospace;white-space:pre-wrap;border-bottom:1px solid var(--bds)}
 .pk{color:#9fb0ff}.pt{color:var(--warn)}.pv{color:#cdd3da}
 #jump{position:absolute;bottom:14px;left:50%;transform:translateX(-50%);background:var(--acc);color:#fff;border:none;display:none;box-shadow:0 4px 16px rgba(0,0,0,.4)}
</style></head><body>
<header>
 <h1>Packet Inspector</h1>
 <span class="stat" id="counts"></span>
 <span class="stat" id="server"></span>
 <div class="spacer"></div>
 <input type="text" id="search" placeholder="filter code or name…"/>
 <label class="chk"><input type="checkbox" id="fe" checked> ev</label>
 <label class="chk"><input type="checkbox" id="fr" checked> req</label>
 <label class="chk"><input type="checkbox" id="fp" checked> resp</label>
 <label class="chk"><input type="checkbox" id="fu"> unknown</label>
 <label class="chk"><input type="checkbox" id="fn" checked> hide noise</label>
 <button id="mark" class="mark" title="Insert a marker (press M) right when you do an in-game action">⌖ Mark</button>
 <button id="pause">Pause</button>
 <button id="clear">Clear</button>
</header>
<div class="wrap" id="wrap"><table>
 <thead><tr><th>#</th><th>Time</th><th>Type</th><th>Code</th><th>Name</th><th>rc</th><th>params</th></tr></thead>
 <tbody id="rows"></tbody>
</table><button id="jump">▼ New packets — jump to live</button></div>
<script>
let since=0, paused=false, following=true, openSeq=null;
const all=[];                       // every record
const NOISE_EV=new Set([3,6,7,8,9,12,19,20,21,22]); // Move/Health/Energy/Cast spam (event codes)
const NOISE_OP=new Set([22]);                       // Move (operation code) — request + its response
const $=id=>document.getElementById(id);
const wrap=$('wrap'), rows=$('rows');

function isNoise(r){ return r.type==='event' ? NOISE_EV.has(r.code) : NOISE_OP.has(r.code); }
function visible(r){
  if($('fn').checked && isNoise(r)) return false;
  if(r.type==='event'&&!$('fe').checked) return false;
  if(r.type==='request'&&!$('fr').checked) return false;
  if(r.type==='response'&&!$('fp').checked) return false;
  if($('fu').checked&&r.known) return false;
  const q=$('search').value.trim().toLowerCase();
  if(q&&!(String(r.code).includes(q)||r.name.toLowerCase().includes(q))) return false;
  return true;
}
function esc(x){return String(x).replace(/[&<>]/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;'}[c]))}
function fmtVal(o,ind){
  ind=ind||''; if(o==null) return '<span class="pt">null</span>';
  const t=o.t,v=o.v;
  if(Array.isArray(v)){ let s='<span class="pt">'+esc(t)+'</span>'; v.forEach((e,i)=>{s+='\n'+ind+'  <span class="pk">'+i+':</span> '+fmtVal(e,ind+'  ');}); return s; }
  if(v&&typeof v==='object'&&t==='dict'){ let s='<span class="pt">dict</span>'; for(const k of Object.keys(v)){s+='\n'+ind+'  <span class="pk">'+esc(k)+':</span> '+fmtVal(v[k],ind+'  ');} return s; }
  if(v&&typeof v==='object'){ return '<span class="pt">'+esc(t)+'</span> <span class="pv">'+esc(JSON.stringify(v))+'</span>'; }
  return '<span class="pt">'+esc(t)+'</span> <span class="pv">'+esc(v)+'</span>';
}
function fmtParams(p){ const ks=Object.keys(p).sort((a,b)=>(+a)-(+b)); let s=''; for(const k of ks) s+='<span class="pk">['+k+']</span> '+fmtVal(p[k])+'\n'; return s||'(no params)'; }
function rowEl(r){
  const tr=document.createElement('tr'); tr.className='r';
  tr.innerHTML='<td class="mono dim">'+r.seq+'</td><td class="mono dim">'+r.time+'</td>'+
    '<td><span class="badge '+r.type+'">'+r.type+'</span></td><td class="mono">'+r.code+'</td>'+
    '<td class="'+(r.known?'':'unk')+'">'+esc(r.name)+(r.known?'':' ⚠')+'</td>'+
    '<td class="mono dim">'+(r.type==='response'?r.returnCode:'')+'</td><td class="mono dim">'+r.numParams+' ▸</td>';
  tr.onclick=()=>{ openSeq=(openSeq===r.seq?null:r.seq); rebuild(); };
  return tr;
}
function append(r){
  if(!visible(r)) return;
  rows.appendChild(rowEl(r));
  if(openSeq===r.seq){ const d=document.createElement('tr'); d.innerHTML='<td colspan="7" style="padding:0"><pre>'+fmtParams(r.params)+'</pre></td>'; rows.appendChild(d); }
}
function rebuild(){
  rows.innerHTML='';
  const list=all.filter(visible).slice(-3000);
  for(const r of list) append(r);
  if(following) wrap.scrollTop=wrap.scrollHeight;
}
function addMark(){
  const t=new Date().toLocaleTimeString();
  const tr=document.createElement('tr'); tr.className='mark';
  tr.innerHTML='<td colspan="7">— MARK '+t+' —</td>'; rows.appendChild(tr);
  if(following) wrap.scrollTop=wrap.scrollHeight;
}
async function poll(){
  if(!paused){ try{
    const d=await (await fetch('/api/recent?since='+since)).json();
    since=d.seq;
    if(d.records&&d.records.length){ for(const r of d.records){ all.push(r); append(r); } if(following) wrap.scrollTop=wrap.scrollHeight; }
    $('counts').innerHTML='<b>'+(d.counts.event||0)+'</b> ev · <b>'+(d.counts.request||0)+'</b> req · <b>'+(d.counts.response||0)+'</b> resp · <b class="unkc">'+(d.counts.unknown||0)+'</b> unknown';
    $('server').textContent=d.server?('server: '+d.server):'';
  }catch(e){} }
  setTimeout(poll,500);
}
// stop-on-scroll: leaving the bottom pauses auto-follow; scroll back down (or click the pill) to resume.
wrap.addEventListener('scroll',()=>{
  const atBottom = wrap.scrollHeight - wrap.scrollTop - wrap.clientHeight < 40;
  following = atBottom;
  $('jump').style.display = atBottom ? 'none' : 'block';
});
$('jump').onclick=()=>{ following=true; $('jump').style.display='none'; wrap.scrollTop=wrap.scrollHeight; };
['fe','fr','fp','fu','fn','search'].forEach(id=>$(id).addEventListener('input',rebuild));
$('pause').onclick=()=>{ paused=!paused; $('pause').textContent=paused?'Resume':'Pause'; $('pause').classList.toggle('on',paused); };
$('clear').onclick=()=>{ all.length=0; rows.innerHTML=''; };
$('mark').onclick=addMark;
document.addEventListener('keydown',e=>{ if(e.key==='m'||e.key==='M'){ if(document.activeElement.tagName!=='INPUT') addMark(); } });
poll();
</script></body></html>`
