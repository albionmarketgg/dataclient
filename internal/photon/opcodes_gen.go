// Code generated. DO NOT EDIT. Albion Photon operation codes (param 253).
package photon

// OperationCode is the Photon operation code (param 253). Values are implicit ordinals.
type OperationCode int16

const (
	OpUnused OperationCode = iota
	OpPing
	OpJoin
	OpVersionedOperation
	OpCreateAccount
	OpLogin
	OpCreateGuestAccount
	OpCreatePlatformOnlyAccount
	OpSendCrashLog
	OpSendTraceRoute
	OpSendVfxStats
	OpSendGamePingInfo
	OpCreateCharacter
	OpDeleteCharacter
	OpSelectCharacter
	OpAcceptPopups
	OpRedeemKeycode
	OpGetGameServerByCluster
	OpGetShopPurchaseUrl
	OpGetReferralSeasonDetails
	OpGetReferralLink
	OpGetShopTilesForCategory
	OpMove
	OpAttackStart
	OpCastStart
	OpCastCancel
	OpTerminateToggleSpell
	OpChannelingCancel
	OpAttackBuildingStart
	OpInventoryDestroyItem
	OpInventoryMoveItem
	OpInventoryRecoverItem
	OpInventoryRecoverAllItems
	OpInventorySplitStack
	OpInventorySplitStackInto
	OpInventoryStack
	OpInventoryReorder
	OpInventoryDropAll
	OpInventoryAddToStacks
	OpInventoryMoveGivenItems
	OpGetClusterData
	OpChangeCluster
	OpConsoleCommand
	OpChatMessage
	OpReportClientError
	OpRegisterToObject
	OpUnRegisterFromObject
	OpCraftBuildingChangeSettings
	OpCraftBuildingTakeMoney
	OpRepairBuildingChangeSettings
	OpRepairBuildingTakeMoney
	OpActionBuildingChangeSettings
	OpHarvestStart
	OpHarvestCancel
	OpTakeSilver
	OpActionOnBuildingStart
	OpActionOnBuildingCancel
	OpInstallResourceStart
	OpInstallResourceCancel
	OpInstallSilver
	OpBuildingFillNutrition
	OpBuildingChangeRenovationState
	OpBuildingBuySkin
	OpBuildingClaim
	OpBuildingGiveup
	OpBuildingNutritionSilverStorageDeposit
	OpBuildingNutritionSilverStorageWithdraw
	OpBuildingNutritionSilverRewardSet
	OpConstructionSiteCreate
	OpPlaceableObjectPlace
	OpPlaceableObjectPlaceCancel
	OpPlaceableObjectPickup
	OpFurnitureObjectUse
	OpFarmableHarvest
	OpFarmableFinishGrownItem
	OpFarmableDestroy
	OpFarmableGetProduct
	OpFarmableFill
	OpTearDownConstructionSite
	OpAuctionCreateOffer
	OpAuctionCreateRequest
	OpAuctionGetOffers
	OpAuctionGetRequests
	OpAuctionBuyOffer
	OpAuctionAbortAuction
	OpAuctionModifyAuction
	OpAuctionAbortOffer
	OpAuctionAbortRequest
	OpAuctionSellRequest
	OpAuctionGetFinishedAuctions
	OpAuctionGetFinishedAuctionsCount
	OpAuctionFetchAuction
	OpAuctionGetMyOpenOffers
	OpAuctionGetMyOpenRequests
	OpAuctionGetMyOpenAuctions
	OpAuctionGetItemAverageStats
	OpAuctionGetItemAverageValue
	OpAuctionGetLowestOfferPrices
	OpContainerOpen
	OpContainerClose
	OpContainerManageSubContainer
	OpRespawn
	OpSuicide
	OpJoinGuild
	OpLeaveGuild
	OpCreateGuild
	OpInviteToGuild
	OpDeclineGuildInvitation
	OpKickFromGuild
	OpInstantJoinGuild
	OpDuellingChallengePlayer
	OpDuellingAcceptChallenge
	OpDuellingDenyChallenge
	OpChangeClusterTax
	OpClaimTerritory
	OpGiveUpTerritory
	OpChangeTerritoryAccessRights
	OpGetMonolithInfo
	OpGetClaimInfo
	OpGetAttackInfo
	OpGetTerritorySeasonPoints
	OpGetAttackSchedule
	OpGetMatches
	OpGetMatchDetails
	OpJoinMatch
	OpLeaveMatch
	OpGetClusterInstanceInfoForStaticCluster
	OpChangeChatSettings
	OpLogoutStart
	OpLogoutCancel
	OpClaimOrbStart
	OpClaimOrbCancel
	OpMatchLootChestOpeningStart
	OpMatchLootChestOpeningCancel
	OpDepositToGuildAccount
	OpWithdrawalFromAccount
	OpChangeGuildPayUpkeepFlag
	OpChangeGuildTax
	OpGetMyTerritories
	OpMorganaCommand
	OpGetServerInfo
	OpSubscribeToCluster
	OpAnswerMercenaryInvitation
	OpGetCharacterEquipment
	OpGetCharacterSteamAchievements
	OpGetCharacterStats
	OpGetKillHistoryDetails
	OpReSpecAchievement
	OpChangeAvatar
	OpGetRankings
	OpGetRank
	OpGetGvgSeasonRankings
	OpGetGvgSeasonRank
	OpGetGvgSeasonHistoryRankings
	OpGetGvgSeasonGuildMemberHistory
	OpKickFromGvGMatch
	OpGetCrystalLeagueDailySeasonPoints
	OpGetChestLogs
	OpGetAccessRightLogs
	OpGetGuildAccountLogs
	OpGetGuildAccountLogsLargeAmount
	OpInviteToPlayerTrade
	OpPlayerTradeCancel
	OpPlayerTradeInvitationAccept
	OpPlayerTradeAddItem
	OpPlayerTradeRemoveItem
	OpPlayerTradeAcceptTrade
	OpPlayerTradeSetSilverOrGold
	OpSendMiniMapPing
	OpStuck
	OpBuyRealEstate
	OpClaimRealEstate
	OpGiveUpRealEstate
	OpChangeRealEstateOutline
	OpGetMailInfos
	OpGetMailCount
	OpReadMail
	OpSendNewMail
	OpDeleteMail
	OpMarkMailUnread
	OpClaimAttachmentFromMail
	OpApplyToGuild
	OpAnswerGuildApplication
	OpRequestGuildFinderFilteredList
	OpUpdateGuildRecruitmentInfo
	OpRequestGuildRecruitmentInfo
	OpRequestGuildFinderNameSearch
	OpRequestGuildFinderRecommendedList
	OpRegisterChatPeer
	OpSendChatMessage
	OpSendModeratorMessage
	OpJoinChatChannel
	OpLeaveChatChannel
	OpSendWhisperMessage
	OpSay
	OpPlayEmote
	OpStopEmote
	OpGetClusterMapInfo
	OpAccessRightsChangeSettings
	OpMount
	OpMountCancel
	OpBuyJourney
	OpSetSaleStatusForEstate
	OpResolveGuildOrPlayerName
	OpGetRespawnInfos
	OpMakeHome
	OpLeaveHome
	OpResurrectionReply
	OpAllianceCreate
	OpAllianceDisband
	OpAllianceGetMemberInfos
	OpAllianceInvite
	OpAllianceAnswerInvitation
	OpAllianceCancelInvitation
	OpAllianceKickGuild
	OpAllianceLeave
	OpAllianceChangeGoldPaymentFlag
	OpAllianceGetDetailInfo
	OpGetIslandInfos
	OpBuyMyIsland
	OpBuyGuildIsland
	OpUpgradeMyIsland
	OpUpgradeGuildIsland
	OpTerritoryFillNutrition
	OpTeleportBack
	OpPartyInvitePlayer
	OpPartyRequestJoin
	OpPartyAnswerInvitation
	OpPartyAnswerJoinRequest
	OpPartyLeave
	OpPartyKickPlayer
	OpPartyMakeLeader
	OpPartyChangeLootSetting
	OpPartyMarkObject
	OpPartySetRole
	OpPartyChangeFactionWarfareRequestReinforcementsSetting
	OpSetGuildCodex
	OpExitEnterStart
	OpExitEnterCancel
	OpQuestGiverRequest
	OpGoldMarketGetBuyOffer
	OpGoldMarketGetBuyOfferFromSilver
	OpGoldMarketGetSellOffer
	OpGoldMarketGetSellOfferFromSilver
	OpGoldMarketBuyGold
	OpGoldMarketSellGold
	OpGoldMarketCreateSellOrder
	OpGoldMarketCreateBuyOrder
	OpGoldMarketGetInfos
	OpGoldMarketCancelOrder
	OpGoldMarketGetAverageInfo
	OpTreasureChestUsingStart
	OpTreasureChestUsingCancel
	OpUseLootChest
	OpUseShrine
	OpUseHellgateShrine
	OpGetSiegeBannerInfo
	OpLaborerStartJob
	OpLaborerTakeJobLoot
	OpLaborerDismiss
	OpLaborerMove
	OpLaborerBuyItem
	OpLaborerUpgrade
	OpBuyPremium
	OpRealEstateGetAuctionData
	OpRealEstateBidOnAuction
	OpFriendInvite
	OpFriendAnswerInvitation
	OpFriendCancelnvitation
	OpFriendRemove
	OpEquipmentItemChangeSpell
	OpExpeditionRegister
	OpExpeditionRegisterCancel
	OpJoinExpedition
	OpDeclineExpeditionInvitation
	OpVoteStart
	OpVoteDoVote
	OpRatingDoRate
	OpEnteringExpeditionStart
	OpEnteringExpeditionCancel
	OpActivateExpeditionCheckPoint
	OpArenaRegister
	OpArenaAddInvite
	OpArenaRegisterCancel
	OpArenaLeave
	OpJoinArenaMatch
	OpDeclineArenaInvitation
	OpEnteringArenaStart
	OpEnteringArenaCancel
	OpArenaCustomMatch
	OpUpdateCharacterStatement
	OpBoostFarmable
	OpGetStrikeHistory
	OpUseFunction
	OpUsePortalEntrance
	OpResetPortalBinding
	OpQueryPortalBinding
	OpClaimPaymentTransaction
	OpChangeUseFlag
	OpClientPerformanceStats
	OpExtendedHardwareStats
	OpClientLowMemoryWarning
	OpTerritoryClaimStart
	OpTerritoryClaimCancel
	OpDeliverCarriableObjectStart
	OpDeliverCarriableObjectCancel
	OpTerritoryUpgradeWithPowerCrystal
	OpRequestAppStoreProducts
	OpVerifyProductPurchase
	OpQueryGuildPlayerStats
	OpQueryAllianceGuildStats
	OpTrackAchievements
	OpSetAchievementsAutoLearn
	OpDepositItemToGuildCurrency
	OpWithdrawalItemFromGuildCurrency
	OpAuctionSellSpecificItemRequest
	OpFishingStart
	OpFishingCasting
	OpFishingCast
	OpFishingCatch
	OpFishingPull
	OpFishingGiveLine
	OpFishingFinish
	OpFishingCancel
	OpCreateGuildAccessTag
	OpDeleteGuildAccessTag
	OpRenameGuildAccessTag
	OpFlagGuildAccessTagGuildPermission
	OpAssignGuildAccessTag
	OpRemoveGuildAccessTagFromPlayer
	OpModifyGuildAccessTagEditors
	OpRequestPublicAccessTags
	OpChangeAccessTagPublicFlag
	OpUpdateGuildAccessTag
	OpSteamStartMicrotransaction
	OpSteamFinishMicrotransaction
	OpRequestXboxPurchaseIntent
	OpCloseXboxPurchaseIntent
	OpSteamIdHasActiveAccount
	OpCheckEmailAccountState
	OpLinkAccountToSteamId
	OpEpicIdHasActiveAccount
	OpLinkAccountToEpicId
	OpXboxIdHasActiveAccount
	OpInAppConfirmPaymentGooglePlay
	OpInAppConfirmPaymentAppleAppStore
	OpInAppPurchaseRequest
	OpInAppPurchaseFailed
	OpCharacterSubscriptionInfo
	OpAccountSubscriptionInfo
	OpBuyGvgSeasonBooster
	OpChangeFlaggingPrepare
	OpOverCharge
	OpOverChargeEnd
	OpRequestTrusted
	OpChangeGuildLogo
	OpPartyFinderRegisterForUpdates
	OpPartyFinderUnregisterForUpdates
	OpPartyFinderEnlistNewPartySearch
	OpPartyFinderDeletePartySearch
	OpPartyFinderChangePartySearch
	OpPartyFinderChangeRole
	OpPartyFinderApplyForGroup
	OpPartyFinderAcceptOrDeclineApplyForGroup
	OpPartyFinderGetEquipmentSnapshot
	OpPartyFinderRegisterApplicants
	OpPartyFinderUnregisterApplicants
	OpPartyFinderFulltextSearch
	OpPartyFinderRequestEquipmentSnapshot
	OpGetPersonalSeasonTrackerData
	OpGetPersonalSeasonPastRewardData
	OpUseConsumableFromInventory
	OpClaimPersonalSeasonReward
	OpXignCodeMessageToServer
	OpBattlEyeMessageToServer
	OpSetNextTutorialState
	OpAddPlayerToMuteList
	OpRemovePlayerFromMuteList
	OpProductShopUserEvent
	OpGetVanityUnlocks
	OpBuyVanityUnlocks
	OpGetMountSkins
	OpSetMountSkin
	OpSetWardrobe
	OpChangeCustomization
	OpChangePlayerIslandData
	OpGetGuildChallengePoints
	OpSmartQueueJoin
	OpSmartQueueLeave
	OpSmartQueueSelectSpawnCluster
	OpUpgradeHideout
	OpInitHideoutAttackStart
	OpInitHideoutAttackCancel
	OpHideoutFillNutrition
	OpHideoutGetInfo
	OpHideoutGetOwnerInfo
	OpHideoutSetTribute
	OpHideoutUpgradeWithPowerCrystal
	OpHideoutDeclareHQ
	OpHideoutUndeclareHQ
	OpHideoutGetHQRequirements
	OpHideoutBoost
	OpHideoutBoostConstruction
	OpOpenWorldAttackScheduleStart
	OpOpenWorldAttackScheduleCancel
	OpOpenWorldAttackConquerStart
	OpOpenWorldAttackConquerCancel
	OpGetOpenWorldAttackDetails
	OpGetNextOpenWorldAttackScheduleTime
	OpRecoverVaultFromHideout
	OpGetGuildEnergyDrainInfo
	OpChannelingUpdate
	OpUseCorruptedShrine
	OpRequestEstimatedMarketValue
	OpLogFeedback
	OpGetInfamyInfo
	OpGetPartySmartClusterQueuePriority
	OpSetPartySmartClusterQueuePriority
	OpClientAntiAutoClickerInfo
	OpClientBotPatternDetectionInfo
	OpClientAntiGatherClickerInfo
	OpLoadoutCreate
	OpLoadoutRead
	OpLoadoutReadHeaders
	OpLoadoutUpdate
	OpLoadoutDelete
	OpLoadoutOrderUpdate
	OpLoadoutEquip
	OpBatchUseItemCancel
	OpEnlistFactionWarfare
	OpGetFactionWarfareWeeklyReport
	OpClaimFactionWarfareWeeklyReport
	OpGetFactionWarfareCampaignData
	OpClaimFactionWarfareItemReward
	OpSendMemoryConsumption
	OpPickupCarriableObjectStart
	OpPickupCarriableObjectCancel
	OpSetSavingChestLogsFlag
	OpGetSavingChestLogsFlag
	OpRegisterGuestAccount
	OpResendGuestAccountVerificationEmail
	OpDoSimpleActionStart
	OpDoSimpleActionCancel
	OpGetGvgSeasonContributionByActivity
	OpGetGvgSeasonContributionByCrystalLeague
	OpGetGuildMightCategoryContribution
	OpGetGuildMightCategoryOverview
	OpGetPvpChallengeData
	OpClaimPvpChallengeWeeklyReward
	OpGetPersonalMightStats
	OpGetPvpChallengeSeasonRewards
	OpGetPvpChallengeSeasonRewardItems
	OpClaimPvpChallengeSeasonRewards
	OpClaimPvpChallengeSeasonRewardItems
	OpAuctionGetLoadoutOffers
	OpAuctionBuyLoadoutOffer
	OpAccountDeletionRequest
	OpAccountReactivationRequest
	OpCreateModeratorNotesForAccount
	OpGetModeratorNotesForAccount
	OpGetModerationEscalationDefiniton
	OpEventBasedPopupAddSeen
	OpGetItemKillHistory
	OpGetVanityConsumables
	OpEquipKillEmote
	OpChangeKillEmotePlayOnKnockdownSetting
	OpBuyVanityConsumableCharges
	OpReclaimVanityItem
	OpGetArenaRankings
	OpGetCrystalLeagueStatistics
	OpSendOptionsLog
	OpSendControlsOptionsLog
	OpMistsUseImmediateReturnExit
	OpMistsUseStaticEntrance
	OpMistsUseCityRoadsEntrance
	OpChangeNewGuildMemberMail
	OpGetNewGuildMemberMail
	OpChangeGuildFactionAllegiance
	OpGetGuildFactionAllegiance
	OpGuildBannerChange
	OpGuildGetOptionalStats
	OpGuildSetOptionalStats
	OpGetPlayerInfoForStalk
	OpPayGoldForCharacterTypeChange
	OpQuickSellAuctionQueryAction
	OpQuickSellAuctionSellAction
	OpFcmTokenToServer
	OpApnsTokenToServer
	OpDeathRecap
	OpAuctionFetchFinishedAuctions
	OpAbortAuctionFetchFinishedAuctions
	OpRequestLegendaryEvenHistory
	OpPartyAnswerStartHuntRequest
	OpHuntAbort
	OpUseFindTrackSpellFromItemPrepare
	OpInteractWithTrackStart
	OpInteractWithTrackCancel
	OpTerritoryRaidStart
	OpTerritoryRaidCancel
	OpTerritoryClaimRaidedRawEnergyCrystalResult
	OpGvGSeasonPlayerGuildParticipationDetails
	OpDailyMightBonus
	OpClaimDailyMightBonus
	OpGetFortificationGroupInfo
	OpUpgradeFortificationGroup
	OpCancelUpgradeFortificationGroup
	OpDowngradeFortificationGroup
	OpGetClusterActivityChestEstimates
	OpPartyReadyCheckBegin
	OpPartyReadyCheckUpdate
	OpClaimAlbionJournalReward
	OpTrackAlbionJournalAchievements
	OpTrackAlbionJournalAchievementSubCategory
	OpRequestOutlandsTeleportationUsage
	OpPickupFromPiledObjectStart
	OpPickupFromPiledObjectCancel
	OpAssetOverview
	OpAssetOverviewTabs
	OpAssetOverviewTabContent
	OpAssetOverviewUnfreezeCache
	OpAssetOverviewSearch
	OpAssetOverviewSearchTabs
	OpAssetOverviewSearchTabContent
	OpAssetOverviewRecoverPlayerVault
	OpImmortalizeKillTrophy
	OpArmorySearch
	OpArmoryItemUsageStatistics
	OpArmoryActivityUsageStatistics
	OpHellDungeonUseStaticEntrance
	OpTravelIslandShowroom
	OpGetXuids
	OpXboxServiceTicket
	OpEvaluatePlatformPerks
	OpLinkAccountToXbox
	OpTravelFactionWarfarePortal
	OpRequestRedZoneEventStandings
	OpGetZergDebuffInfo
	OpRequestLoreSnippetStates
	OpRetrieveCarriableObjectStart
	OpRetrieveCarriableObjectCancel
)

var operationCodeNames = map[OperationCode]string{
	OpUnused: "Unused",
	OpPing: "Ping",
	OpJoin: "Join",
	OpVersionedOperation: "VersionedOperation",
	OpCreateAccount: "CreateAccount",
	OpLogin: "Login",
	OpCreateGuestAccount: "CreateGuestAccount",
	OpCreatePlatformOnlyAccount: "CreatePlatformOnlyAccount",
	OpSendCrashLog: "SendCrashLog",
	OpSendTraceRoute: "SendTraceRoute",
	OpSendVfxStats: "SendVfxStats",
	OpSendGamePingInfo: "SendGamePingInfo",
	OpCreateCharacter: "CreateCharacter",
	OpDeleteCharacter: "DeleteCharacter",
	OpSelectCharacter: "SelectCharacter",
	OpAcceptPopups: "AcceptPopups",
	OpRedeemKeycode: "RedeemKeycode",
	OpGetGameServerByCluster: "GetGameServerByCluster",
	OpGetShopPurchaseUrl: "GetShopPurchaseUrl",
	OpGetReferralSeasonDetails: "GetReferralSeasonDetails",
	OpGetReferralLink: "GetReferralLink",
	OpGetShopTilesForCategory: "GetShopTilesForCategory",
	OpMove: "Move",
	OpAttackStart: "AttackStart",
	OpCastStart: "CastStart",
	OpCastCancel: "CastCancel",
	OpTerminateToggleSpell: "TerminateToggleSpell",
	OpChannelingCancel: "ChannelingCancel",
	OpAttackBuildingStart: "AttackBuildingStart",
	OpInventoryDestroyItem: "InventoryDestroyItem",
	OpInventoryMoveItem: "InventoryMoveItem",
	OpInventoryRecoverItem: "InventoryRecoverItem",
	OpInventoryRecoverAllItems: "InventoryRecoverAllItems",
	OpInventorySplitStack: "InventorySplitStack",
	OpInventorySplitStackInto: "InventorySplitStackInto",
	OpInventoryStack: "InventoryStack",
	OpInventoryReorder: "InventoryReorder",
	OpInventoryDropAll: "InventoryDropAll",
	OpInventoryAddToStacks: "InventoryAddToStacks",
	OpInventoryMoveGivenItems: "InventoryMoveGivenItems",
	OpGetClusterData: "GetClusterData",
	OpChangeCluster: "ChangeCluster",
	OpConsoleCommand: "ConsoleCommand",
	OpChatMessage: "ChatMessage",
	OpReportClientError: "ReportClientError",
	OpRegisterToObject: "RegisterToObject",
	OpUnRegisterFromObject: "UnRegisterFromObject",
	OpCraftBuildingChangeSettings: "CraftBuildingChangeSettings",
	OpCraftBuildingTakeMoney: "CraftBuildingTakeMoney",
	OpRepairBuildingChangeSettings: "RepairBuildingChangeSettings",
	OpRepairBuildingTakeMoney: "RepairBuildingTakeMoney",
	OpActionBuildingChangeSettings: "ActionBuildingChangeSettings",
	OpHarvestStart: "HarvestStart",
	OpHarvestCancel: "HarvestCancel",
	OpTakeSilver: "TakeSilver",
	OpActionOnBuildingStart: "ActionOnBuildingStart",
	OpActionOnBuildingCancel: "ActionOnBuildingCancel",
	OpInstallResourceStart: "InstallResourceStart",
	OpInstallResourceCancel: "InstallResourceCancel",
	OpInstallSilver: "InstallSilver",
	OpBuildingFillNutrition: "BuildingFillNutrition",
	OpBuildingChangeRenovationState: "BuildingChangeRenovationState",
	OpBuildingBuySkin: "BuildingBuySkin",
	OpBuildingClaim: "BuildingClaim",
	OpBuildingGiveup: "BuildingGiveup",
	OpBuildingNutritionSilverStorageDeposit: "BuildingNutritionSilverStorageDeposit",
	OpBuildingNutritionSilverStorageWithdraw: "BuildingNutritionSilverStorageWithdraw",
	OpBuildingNutritionSilverRewardSet: "BuildingNutritionSilverRewardSet",
	OpConstructionSiteCreate: "ConstructionSiteCreate",
	OpPlaceableObjectPlace: "PlaceableObjectPlace",
	OpPlaceableObjectPlaceCancel: "PlaceableObjectPlaceCancel",
	OpPlaceableObjectPickup: "PlaceableObjectPickup",
	OpFurnitureObjectUse: "FurnitureObjectUse",
	OpFarmableHarvest: "FarmableHarvest",
	OpFarmableFinishGrownItem: "FarmableFinishGrownItem",
	OpFarmableDestroy: "FarmableDestroy",
	OpFarmableGetProduct: "FarmableGetProduct",
	OpFarmableFill: "FarmableFill",
	OpTearDownConstructionSite: "TearDownConstructionSite",
	OpAuctionCreateOffer: "AuctionCreateOffer",
	OpAuctionCreateRequest: "AuctionCreateRequest",
	OpAuctionGetOffers: "AuctionGetOffers",
	OpAuctionGetRequests: "AuctionGetRequests",
	OpAuctionBuyOffer: "AuctionBuyOffer",
	OpAuctionAbortAuction: "AuctionAbortAuction",
	OpAuctionModifyAuction: "AuctionModifyAuction",
	OpAuctionAbortOffer: "AuctionAbortOffer",
	OpAuctionAbortRequest: "AuctionAbortRequest",
	OpAuctionSellRequest: "AuctionSellRequest",
	OpAuctionGetFinishedAuctions: "AuctionGetFinishedAuctions",
	OpAuctionGetFinishedAuctionsCount: "AuctionGetFinishedAuctionsCount",
	OpAuctionFetchAuction: "AuctionFetchAuction",
	OpAuctionGetMyOpenOffers: "AuctionGetMyOpenOffers",
	OpAuctionGetMyOpenRequests: "AuctionGetMyOpenRequests",
	OpAuctionGetMyOpenAuctions: "AuctionGetMyOpenAuctions",
	OpAuctionGetItemAverageStats: "AuctionGetItemAverageStats",
	OpAuctionGetItemAverageValue: "AuctionGetItemAverageValue",
	OpAuctionGetLowestOfferPrices: "AuctionGetLowestOfferPrices",
	OpContainerOpen: "ContainerOpen",
	OpContainerClose: "ContainerClose",
	OpContainerManageSubContainer: "ContainerManageSubContainer",
	OpRespawn: "Respawn",
	OpSuicide: "Suicide",
	OpJoinGuild: "JoinGuild",
	OpLeaveGuild: "LeaveGuild",
	OpCreateGuild: "CreateGuild",
	OpInviteToGuild: "InviteToGuild",
	OpDeclineGuildInvitation: "DeclineGuildInvitation",
	OpKickFromGuild: "KickFromGuild",
	OpInstantJoinGuild: "InstantJoinGuild",
	OpDuellingChallengePlayer: "DuellingChallengePlayer",
	OpDuellingAcceptChallenge: "DuellingAcceptChallenge",
	OpDuellingDenyChallenge: "DuellingDenyChallenge",
	OpChangeClusterTax: "ChangeClusterTax",
	OpClaimTerritory: "ClaimTerritory",
	OpGiveUpTerritory: "GiveUpTerritory",
	OpChangeTerritoryAccessRights: "ChangeTerritoryAccessRights",
	OpGetMonolithInfo: "GetMonolithInfo",
	OpGetClaimInfo: "GetClaimInfo",
	OpGetAttackInfo: "GetAttackInfo",
	OpGetTerritorySeasonPoints: "GetTerritorySeasonPoints",
	OpGetAttackSchedule: "GetAttackSchedule",
	OpGetMatches: "GetMatches",
	OpGetMatchDetails: "GetMatchDetails",
	OpJoinMatch: "JoinMatch",
	OpLeaveMatch: "LeaveMatch",
	OpGetClusterInstanceInfoForStaticCluster: "GetClusterInstanceInfoForStaticCluster",
	OpChangeChatSettings: "ChangeChatSettings",
	OpLogoutStart: "LogoutStart",
	OpLogoutCancel: "LogoutCancel",
	OpClaimOrbStart: "ClaimOrbStart",
	OpClaimOrbCancel: "ClaimOrbCancel",
	OpMatchLootChestOpeningStart: "MatchLootChestOpeningStart",
	OpMatchLootChestOpeningCancel: "MatchLootChestOpeningCancel",
	OpDepositToGuildAccount: "DepositToGuildAccount",
	OpWithdrawalFromAccount: "WithdrawalFromAccount",
	OpChangeGuildPayUpkeepFlag: "ChangeGuildPayUpkeepFlag",
	OpChangeGuildTax: "ChangeGuildTax",
	OpGetMyTerritories: "GetMyTerritories",
	OpMorganaCommand: "MorganaCommand",
	OpGetServerInfo: "GetServerInfo",
	OpSubscribeToCluster: "SubscribeToCluster",
	OpAnswerMercenaryInvitation: "AnswerMercenaryInvitation",
	OpGetCharacterEquipment: "GetCharacterEquipment",
	OpGetCharacterSteamAchievements: "GetCharacterSteamAchievements",
	OpGetCharacterStats: "GetCharacterStats",
	OpGetKillHistoryDetails: "GetKillHistoryDetails",
	OpReSpecAchievement: "ReSpecAchievement",
	OpChangeAvatar: "ChangeAvatar",
	OpGetRankings: "GetRankings",
	OpGetRank: "GetRank",
	OpGetGvgSeasonRankings: "GetGvgSeasonRankings",
	OpGetGvgSeasonRank: "GetGvgSeasonRank",
	OpGetGvgSeasonHistoryRankings: "GetGvgSeasonHistoryRankings",
	OpGetGvgSeasonGuildMemberHistory: "GetGvgSeasonGuildMemberHistory",
	OpKickFromGvGMatch: "KickFromGvGMatch",
	OpGetCrystalLeagueDailySeasonPoints: "GetCrystalLeagueDailySeasonPoints",
	OpGetChestLogs: "GetChestLogs",
	OpGetAccessRightLogs: "GetAccessRightLogs",
	OpGetGuildAccountLogs: "GetGuildAccountLogs",
	OpGetGuildAccountLogsLargeAmount: "GetGuildAccountLogsLargeAmount",
	OpInviteToPlayerTrade: "InviteToPlayerTrade",
	OpPlayerTradeCancel: "PlayerTradeCancel",
	OpPlayerTradeInvitationAccept: "PlayerTradeInvitationAccept",
	OpPlayerTradeAddItem: "PlayerTradeAddItem",
	OpPlayerTradeRemoveItem: "PlayerTradeRemoveItem",
	OpPlayerTradeAcceptTrade: "PlayerTradeAcceptTrade",
	OpPlayerTradeSetSilverOrGold: "PlayerTradeSetSilverOrGold",
	OpSendMiniMapPing: "SendMiniMapPing",
	OpStuck: "Stuck",
	OpBuyRealEstate: "BuyRealEstate",
	OpClaimRealEstate: "ClaimRealEstate",
	OpGiveUpRealEstate: "GiveUpRealEstate",
	OpChangeRealEstateOutline: "ChangeRealEstateOutline",
	OpGetMailInfos: "GetMailInfos",
	OpGetMailCount: "GetMailCount",
	OpReadMail: "ReadMail",
	OpSendNewMail: "SendNewMail",
	OpDeleteMail: "DeleteMail",
	OpMarkMailUnread: "MarkMailUnread",
	OpClaimAttachmentFromMail: "ClaimAttachmentFromMail",
	OpApplyToGuild: "ApplyToGuild",
	OpAnswerGuildApplication: "AnswerGuildApplication",
	OpRequestGuildFinderFilteredList: "RequestGuildFinderFilteredList",
	OpUpdateGuildRecruitmentInfo: "UpdateGuildRecruitmentInfo",
	OpRequestGuildRecruitmentInfo: "RequestGuildRecruitmentInfo",
	OpRequestGuildFinderNameSearch: "RequestGuildFinderNameSearch",
	OpRequestGuildFinderRecommendedList: "RequestGuildFinderRecommendedList",
	OpRegisterChatPeer: "RegisterChatPeer",
	OpSendChatMessage: "SendChatMessage",
	OpSendModeratorMessage: "SendModeratorMessage",
	OpJoinChatChannel: "JoinChatChannel",
	OpLeaveChatChannel: "LeaveChatChannel",
	OpSendWhisperMessage: "SendWhisperMessage",
	OpSay: "Say",
	OpPlayEmote: "PlayEmote",
	OpStopEmote: "StopEmote",
	OpGetClusterMapInfo: "GetClusterMapInfo",
	OpAccessRightsChangeSettings: "AccessRightsChangeSettings",
	OpMount: "Mount",
	OpMountCancel: "MountCancel",
	OpBuyJourney: "BuyJourney",
	OpSetSaleStatusForEstate: "SetSaleStatusForEstate",
	OpResolveGuildOrPlayerName: "ResolveGuildOrPlayerName",
	OpGetRespawnInfos: "GetRespawnInfos",
	OpMakeHome: "MakeHome",
	OpLeaveHome: "LeaveHome",
	OpResurrectionReply: "ResurrectionReply",
	OpAllianceCreate: "AllianceCreate",
	OpAllianceDisband: "AllianceDisband",
	OpAllianceGetMemberInfos: "AllianceGetMemberInfos",
	OpAllianceInvite: "AllianceInvite",
	OpAllianceAnswerInvitation: "AllianceAnswerInvitation",
	OpAllianceCancelInvitation: "AllianceCancelInvitation",
	OpAllianceKickGuild: "AllianceKickGuild",
	OpAllianceLeave: "AllianceLeave",
	OpAllianceChangeGoldPaymentFlag: "AllianceChangeGoldPaymentFlag",
	OpAllianceGetDetailInfo: "AllianceGetDetailInfo",
	OpGetIslandInfos: "GetIslandInfos",
	OpBuyMyIsland: "BuyMyIsland",
	OpBuyGuildIsland: "BuyGuildIsland",
	OpUpgradeMyIsland: "UpgradeMyIsland",
	OpUpgradeGuildIsland: "UpgradeGuildIsland",
	OpTerritoryFillNutrition: "TerritoryFillNutrition",
	OpTeleportBack: "TeleportBack",
	OpPartyInvitePlayer: "PartyInvitePlayer",
	OpPartyRequestJoin: "PartyRequestJoin",
	OpPartyAnswerInvitation: "PartyAnswerInvitation",
	OpPartyAnswerJoinRequest: "PartyAnswerJoinRequest",
	OpPartyLeave: "PartyLeave",
	OpPartyKickPlayer: "PartyKickPlayer",
	OpPartyMakeLeader: "PartyMakeLeader",
	OpPartyChangeLootSetting: "PartyChangeLootSetting",
	OpPartyMarkObject: "PartyMarkObject",
	OpPartySetRole: "PartySetRole",
	OpPartyChangeFactionWarfareRequestReinforcementsSetting: "PartyChangeFactionWarfareRequestReinforcementsSetting",
	OpSetGuildCodex: "SetGuildCodex",
	OpExitEnterStart: "ExitEnterStart",
	OpExitEnterCancel: "ExitEnterCancel",
	OpQuestGiverRequest: "QuestGiverRequest",
	OpGoldMarketGetBuyOffer: "GoldMarketGetBuyOffer",
	OpGoldMarketGetBuyOfferFromSilver: "GoldMarketGetBuyOfferFromSilver",
	OpGoldMarketGetSellOffer: "GoldMarketGetSellOffer",
	OpGoldMarketGetSellOfferFromSilver: "GoldMarketGetSellOfferFromSilver",
	OpGoldMarketBuyGold: "GoldMarketBuyGold",
	OpGoldMarketSellGold: "GoldMarketSellGold",
	OpGoldMarketCreateSellOrder: "GoldMarketCreateSellOrder",
	OpGoldMarketCreateBuyOrder: "GoldMarketCreateBuyOrder",
	OpGoldMarketGetInfos: "GoldMarketGetInfos",
	OpGoldMarketCancelOrder: "GoldMarketCancelOrder",
	OpGoldMarketGetAverageInfo: "GoldMarketGetAverageInfo",
	OpTreasureChestUsingStart: "TreasureChestUsingStart",
	OpTreasureChestUsingCancel: "TreasureChestUsingCancel",
	OpUseLootChest: "UseLootChest",
	OpUseShrine: "UseShrine",
	OpUseHellgateShrine: "UseHellgateShrine",
	OpGetSiegeBannerInfo: "GetSiegeBannerInfo",
	OpLaborerStartJob: "LaborerStartJob",
	OpLaborerTakeJobLoot: "LaborerTakeJobLoot",
	OpLaborerDismiss: "LaborerDismiss",
	OpLaborerMove: "LaborerMove",
	OpLaborerBuyItem: "LaborerBuyItem",
	OpLaborerUpgrade: "LaborerUpgrade",
	OpBuyPremium: "BuyPremium",
	OpRealEstateGetAuctionData: "RealEstateGetAuctionData",
	OpRealEstateBidOnAuction: "RealEstateBidOnAuction",
	OpFriendInvite: "FriendInvite",
	OpFriendAnswerInvitation: "FriendAnswerInvitation",
	OpFriendCancelnvitation: "FriendCancelnvitation",
	OpFriendRemove: "FriendRemove",
	OpEquipmentItemChangeSpell: "EquipmentItemChangeSpell",
	OpExpeditionRegister: "ExpeditionRegister",
	OpExpeditionRegisterCancel: "ExpeditionRegisterCancel",
	OpJoinExpedition: "JoinExpedition",
	OpDeclineExpeditionInvitation: "DeclineExpeditionInvitation",
	OpVoteStart: "VoteStart",
	OpVoteDoVote: "VoteDoVote",
	OpRatingDoRate: "RatingDoRate",
	OpEnteringExpeditionStart: "EnteringExpeditionStart",
	OpEnteringExpeditionCancel: "EnteringExpeditionCancel",
	OpActivateExpeditionCheckPoint: "ActivateExpeditionCheckPoint",
	OpArenaRegister: "ArenaRegister",
	OpArenaAddInvite: "ArenaAddInvite",
	OpArenaRegisterCancel: "ArenaRegisterCancel",
	OpArenaLeave: "ArenaLeave",
	OpJoinArenaMatch: "JoinArenaMatch",
	OpDeclineArenaInvitation: "DeclineArenaInvitation",
	OpEnteringArenaStart: "EnteringArenaStart",
	OpEnteringArenaCancel: "EnteringArenaCancel",
	OpArenaCustomMatch: "ArenaCustomMatch",
	OpUpdateCharacterStatement: "UpdateCharacterStatement",
	OpBoostFarmable: "BoostFarmable",
	OpGetStrikeHistory: "GetStrikeHistory",
	OpUseFunction: "UseFunction",
	OpUsePortalEntrance: "UsePortalEntrance",
	OpResetPortalBinding: "ResetPortalBinding",
	OpQueryPortalBinding: "QueryPortalBinding",
	OpClaimPaymentTransaction: "ClaimPaymentTransaction",
	OpChangeUseFlag: "ChangeUseFlag",
	OpClientPerformanceStats: "ClientPerformanceStats",
	OpExtendedHardwareStats: "ExtendedHardwareStats",
	OpClientLowMemoryWarning: "ClientLowMemoryWarning",
	OpTerritoryClaimStart: "TerritoryClaimStart",
	OpTerritoryClaimCancel: "TerritoryClaimCancel",
	OpDeliverCarriableObjectStart: "DeliverCarriableObjectStart",
	OpDeliverCarriableObjectCancel: "DeliverCarriableObjectCancel",
	OpTerritoryUpgradeWithPowerCrystal: "TerritoryUpgradeWithPowerCrystal",
	OpRequestAppStoreProducts: "RequestAppStoreProducts",
	OpVerifyProductPurchase: "VerifyProductPurchase",
	OpQueryGuildPlayerStats: "QueryGuildPlayerStats",
	OpQueryAllianceGuildStats: "QueryAllianceGuildStats",
	OpTrackAchievements: "TrackAchievements",
	OpSetAchievementsAutoLearn: "SetAchievementsAutoLearn",
	OpDepositItemToGuildCurrency: "DepositItemToGuildCurrency",
	OpWithdrawalItemFromGuildCurrency: "WithdrawalItemFromGuildCurrency",
	OpAuctionSellSpecificItemRequest: "AuctionSellSpecificItemRequest",
	OpFishingStart: "FishingStart",
	OpFishingCasting: "FishingCasting",
	OpFishingCast: "FishingCast",
	OpFishingCatch: "FishingCatch",
	OpFishingPull: "FishingPull",
	OpFishingGiveLine: "FishingGiveLine",
	OpFishingFinish: "FishingFinish",
	OpFishingCancel: "FishingCancel",
	OpCreateGuildAccessTag: "CreateGuildAccessTag",
	OpDeleteGuildAccessTag: "DeleteGuildAccessTag",
	OpRenameGuildAccessTag: "RenameGuildAccessTag",
	OpFlagGuildAccessTagGuildPermission: "FlagGuildAccessTagGuildPermission",
	OpAssignGuildAccessTag: "AssignGuildAccessTag",
	OpRemoveGuildAccessTagFromPlayer: "RemoveGuildAccessTagFromPlayer",
	OpModifyGuildAccessTagEditors: "ModifyGuildAccessTagEditors",
	OpRequestPublicAccessTags: "RequestPublicAccessTags",
	OpChangeAccessTagPublicFlag: "ChangeAccessTagPublicFlag",
	OpUpdateGuildAccessTag: "UpdateGuildAccessTag",
	OpSteamStartMicrotransaction: "SteamStartMicrotransaction",
	OpSteamFinishMicrotransaction: "SteamFinishMicrotransaction",
	OpRequestXboxPurchaseIntent: "RequestXboxPurchaseIntent",
	OpCloseXboxPurchaseIntent: "CloseXboxPurchaseIntent",
	OpSteamIdHasActiveAccount: "SteamIdHasActiveAccount",
	OpCheckEmailAccountState: "CheckEmailAccountState",
	OpLinkAccountToSteamId: "LinkAccountToSteamId",
	OpEpicIdHasActiveAccount: "EpicIdHasActiveAccount",
	OpLinkAccountToEpicId: "LinkAccountToEpicId",
	OpXboxIdHasActiveAccount: "XboxIdHasActiveAccount",
	OpInAppConfirmPaymentGooglePlay: "InAppConfirmPaymentGooglePlay",
	OpInAppConfirmPaymentAppleAppStore: "InAppConfirmPaymentAppleAppStore",
	OpInAppPurchaseRequest: "InAppPurchaseRequest",
	OpInAppPurchaseFailed: "InAppPurchaseFailed",
	OpCharacterSubscriptionInfo: "CharacterSubscriptionInfo",
	OpAccountSubscriptionInfo: "AccountSubscriptionInfo",
	OpBuyGvgSeasonBooster: "BuyGvgSeasonBooster",
	OpChangeFlaggingPrepare: "ChangeFlaggingPrepare",
	OpOverCharge: "OverCharge",
	OpOverChargeEnd: "OverChargeEnd",
	OpRequestTrusted: "RequestTrusted",
	OpChangeGuildLogo: "ChangeGuildLogo",
	OpPartyFinderRegisterForUpdates: "PartyFinderRegisterForUpdates",
	OpPartyFinderUnregisterForUpdates: "PartyFinderUnregisterForUpdates",
	OpPartyFinderEnlistNewPartySearch: "PartyFinderEnlistNewPartySearch",
	OpPartyFinderDeletePartySearch: "PartyFinderDeletePartySearch",
	OpPartyFinderChangePartySearch: "PartyFinderChangePartySearch",
	OpPartyFinderChangeRole: "PartyFinderChangeRole",
	OpPartyFinderApplyForGroup: "PartyFinderApplyForGroup",
	OpPartyFinderAcceptOrDeclineApplyForGroup: "PartyFinderAcceptOrDeclineApplyForGroup",
	OpPartyFinderGetEquipmentSnapshot: "PartyFinderGetEquipmentSnapshot",
	OpPartyFinderRegisterApplicants: "PartyFinderRegisterApplicants",
	OpPartyFinderUnregisterApplicants: "PartyFinderUnregisterApplicants",
	OpPartyFinderFulltextSearch: "PartyFinderFulltextSearch",
	OpPartyFinderRequestEquipmentSnapshot: "PartyFinderRequestEquipmentSnapshot",
	OpGetPersonalSeasonTrackerData: "GetPersonalSeasonTrackerData",
	OpGetPersonalSeasonPastRewardData: "GetPersonalSeasonPastRewardData",
	OpUseConsumableFromInventory: "UseConsumableFromInventory",
	OpClaimPersonalSeasonReward: "ClaimPersonalSeasonReward",
	OpXignCodeMessageToServer: "XignCodeMessageToServer",
	OpBattlEyeMessageToServer: "BattlEyeMessageToServer",
	OpSetNextTutorialState: "SetNextTutorialState",
	OpAddPlayerToMuteList: "AddPlayerToMuteList",
	OpRemovePlayerFromMuteList: "RemovePlayerFromMuteList",
	OpProductShopUserEvent: "ProductShopUserEvent",
	OpGetVanityUnlocks: "GetVanityUnlocks",
	OpBuyVanityUnlocks: "BuyVanityUnlocks",
	OpGetMountSkins: "GetMountSkins",
	OpSetMountSkin: "SetMountSkin",
	OpSetWardrobe: "SetWardrobe",
	OpChangeCustomization: "ChangeCustomization",
	OpChangePlayerIslandData: "ChangePlayerIslandData",
	OpGetGuildChallengePoints: "GetGuildChallengePoints",
	OpSmartQueueJoin: "SmartQueueJoin",
	OpSmartQueueLeave: "SmartQueueLeave",
	OpSmartQueueSelectSpawnCluster: "SmartQueueSelectSpawnCluster",
	OpUpgradeHideout: "UpgradeHideout",
	OpInitHideoutAttackStart: "InitHideoutAttackStart",
	OpInitHideoutAttackCancel: "InitHideoutAttackCancel",
	OpHideoutFillNutrition: "HideoutFillNutrition",
	OpHideoutGetInfo: "HideoutGetInfo",
	OpHideoutGetOwnerInfo: "HideoutGetOwnerInfo",
	OpHideoutSetTribute: "HideoutSetTribute",
	OpHideoutUpgradeWithPowerCrystal: "HideoutUpgradeWithPowerCrystal",
	OpHideoutDeclareHQ: "HideoutDeclareHQ",
	OpHideoutUndeclareHQ: "HideoutUndeclareHQ",
	OpHideoutGetHQRequirements: "HideoutGetHQRequirements",
	OpHideoutBoost: "HideoutBoost",
	OpHideoutBoostConstruction: "HideoutBoostConstruction",
	OpOpenWorldAttackScheduleStart: "OpenWorldAttackScheduleStart",
	OpOpenWorldAttackScheduleCancel: "OpenWorldAttackScheduleCancel",
	OpOpenWorldAttackConquerStart: "OpenWorldAttackConquerStart",
	OpOpenWorldAttackConquerCancel: "OpenWorldAttackConquerCancel",
	OpGetOpenWorldAttackDetails: "GetOpenWorldAttackDetails",
	OpGetNextOpenWorldAttackScheduleTime: "GetNextOpenWorldAttackScheduleTime",
	OpRecoverVaultFromHideout: "RecoverVaultFromHideout",
	OpGetGuildEnergyDrainInfo: "GetGuildEnergyDrainInfo",
	OpChannelingUpdate: "ChannelingUpdate",
	OpUseCorruptedShrine: "UseCorruptedShrine",
	OpRequestEstimatedMarketValue: "RequestEstimatedMarketValue",
	OpLogFeedback: "LogFeedback",
	OpGetInfamyInfo: "GetInfamyInfo",
	OpGetPartySmartClusterQueuePriority: "GetPartySmartClusterQueuePriority",
	OpSetPartySmartClusterQueuePriority: "SetPartySmartClusterQueuePriority",
	OpClientAntiAutoClickerInfo: "ClientAntiAutoClickerInfo",
	OpClientBotPatternDetectionInfo: "ClientBotPatternDetectionInfo",
	OpClientAntiGatherClickerInfo: "ClientAntiGatherClickerInfo",
	OpLoadoutCreate: "LoadoutCreate",
	OpLoadoutRead: "LoadoutRead",
	OpLoadoutReadHeaders: "LoadoutReadHeaders",
	OpLoadoutUpdate: "LoadoutUpdate",
	OpLoadoutDelete: "LoadoutDelete",
	OpLoadoutOrderUpdate: "LoadoutOrderUpdate",
	OpLoadoutEquip: "LoadoutEquip",
	OpBatchUseItemCancel: "BatchUseItemCancel",
	OpEnlistFactionWarfare: "EnlistFactionWarfare",
	OpGetFactionWarfareWeeklyReport: "GetFactionWarfareWeeklyReport",
	OpClaimFactionWarfareWeeklyReport: "ClaimFactionWarfareWeeklyReport",
	OpGetFactionWarfareCampaignData: "GetFactionWarfareCampaignData",
	OpClaimFactionWarfareItemReward: "ClaimFactionWarfareItemReward",
	OpSendMemoryConsumption: "SendMemoryConsumption",
	OpPickupCarriableObjectStart: "PickupCarriableObjectStart",
	OpPickupCarriableObjectCancel: "PickupCarriableObjectCancel",
	OpSetSavingChestLogsFlag: "SetSavingChestLogsFlag",
	OpGetSavingChestLogsFlag: "GetSavingChestLogsFlag",
	OpRegisterGuestAccount: "RegisterGuestAccount",
	OpResendGuestAccountVerificationEmail: "ResendGuestAccountVerificationEmail",
	OpDoSimpleActionStart: "DoSimpleActionStart",
	OpDoSimpleActionCancel: "DoSimpleActionCancel",
	OpGetGvgSeasonContributionByActivity: "GetGvgSeasonContributionByActivity",
	OpGetGvgSeasonContributionByCrystalLeague: "GetGvgSeasonContributionByCrystalLeague",
	OpGetGuildMightCategoryContribution: "GetGuildMightCategoryContribution",
	OpGetGuildMightCategoryOverview: "GetGuildMightCategoryOverview",
	OpGetPvpChallengeData: "GetPvpChallengeData",
	OpClaimPvpChallengeWeeklyReward: "ClaimPvpChallengeWeeklyReward",
	OpGetPersonalMightStats: "GetPersonalMightStats",
	OpGetPvpChallengeSeasonRewards: "GetPvpChallengeSeasonRewards",
	OpGetPvpChallengeSeasonRewardItems: "GetPvpChallengeSeasonRewardItems",
	OpClaimPvpChallengeSeasonRewards: "ClaimPvpChallengeSeasonRewards",
	OpClaimPvpChallengeSeasonRewardItems: "ClaimPvpChallengeSeasonRewardItems",
	OpAuctionGetLoadoutOffers: "AuctionGetLoadoutOffers",
	OpAuctionBuyLoadoutOffer: "AuctionBuyLoadoutOffer",
	OpAccountDeletionRequest: "AccountDeletionRequest",
	OpAccountReactivationRequest: "AccountReactivationRequest",
	OpCreateModeratorNotesForAccount: "CreateModeratorNotesForAccount",
	OpGetModeratorNotesForAccount: "GetModeratorNotesForAccount",
	OpGetModerationEscalationDefiniton: "GetModerationEscalationDefiniton",
	OpEventBasedPopupAddSeen: "EventBasedPopupAddSeen",
	OpGetItemKillHistory: "GetItemKillHistory",
	OpGetVanityConsumables: "GetVanityConsumables",
	OpEquipKillEmote: "EquipKillEmote",
	OpChangeKillEmotePlayOnKnockdownSetting: "ChangeKillEmotePlayOnKnockdownSetting",
	OpBuyVanityConsumableCharges: "BuyVanityConsumableCharges",
	OpReclaimVanityItem: "ReclaimVanityItem",
	OpGetArenaRankings: "GetArenaRankings",
	OpGetCrystalLeagueStatistics: "GetCrystalLeagueStatistics",
	OpSendOptionsLog: "SendOptionsLog",
	OpSendControlsOptionsLog: "SendControlsOptionsLog",
	OpMistsUseImmediateReturnExit: "MistsUseImmediateReturnExit",
	OpMistsUseStaticEntrance: "MistsUseStaticEntrance",
	OpMistsUseCityRoadsEntrance: "MistsUseCityRoadsEntrance",
	OpChangeNewGuildMemberMail: "ChangeNewGuildMemberMail",
	OpGetNewGuildMemberMail: "GetNewGuildMemberMail",
	OpChangeGuildFactionAllegiance: "ChangeGuildFactionAllegiance",
	OpGetGuildFactionAllegiance: "GetGuildFactionAllegiance",
	OpGuildBannerChange: "GuildBannerChange",
	OpGuildGetOptionalStats: "GuildGetOptionalStats",
	OpGuildSetOptionalStats: "GuildSetOptionalStats",
	OpGetPlayerInfoForStalk: "GetPlayerInfoForStalk",
	OpPayGoldForCharacterTypeChange: "PayGoldForCharacterTypeChange",
	OpQuickSellAuctionQueryAction: "QuickSellAuctionQueryAction",
	OpQuickSellAuctionSellAction: "QuickSellAuctionSellAction",
	OpFcmTokenToServer: "FcmTokenToServer",
	OpApnsTokenToServer: "ApnsTokenToServer",
	OpDeathRecap: "DeathRecap",
	OpAuctionFetchFinishedAuctions: "AuctionFetchFinishedAuctions",
	OpAbortAuctionFetchFinishedAuctions: "AbortAuctionFetchFinishedAuctions",
	OpRequestLegendaryEvenHistory: "RequestLegendaryEvenHistory",
	OpPartyAnswerStartHuntRequest: "PartyAnswerStartHuntRequest",
	OpHuntAbort: "HuntAbort",
	OpUseFindTrackSpellFromItemPrepare: "UseFindTrackSpellFromItemPrepare",
	OpInteractWithTrackStart: "InteractWithTrackStart",
	OpInteractWithTrackCancel: "InteractWithTrackCancel",
	OpTerritoryRaidStart: "TerritoryRaidStart",
	OpTerritoryRaidCancel: "TerritoryRaidCancel",
	OpTerritoryClaimRaidedRawEnergyCrystalResult: "TerritoryClaimRaidedRawEnergyCrystalResult",
	OpGvGSeasonPlayerGuildParticipationDetails: "GvGSeasonPlayerGuildParticipationDetails",
	OpDailyMightBonus: "DailyMightBonus",
	OpClaimDailyMightBonus: "ClaimDailyMightBonus",
	OpGetFortificationGroupInfo: "GetFortificationGroupInfo",
	OpUpgradeFortificationGroup: "UpgradeFortificationGroup",
	OpCancelUpgradeFortificationGroup: "CancelUpgradeFortificationGroup",
	OpDowngradeFortificationGroup: "DowngradeFortificationGroup",
	OpGetClusterActivityChestEstimates: "GetClusterActivityChestEstimates",
	OpPartyReadyCheckBegin: "PartyReadyCheckBegin",
	OpPartyReadyCheckUpdate: "PartyReadyCheckUpdate",
	OpClaimAlbionJournalReward: "ClaimAlbionJournalReward",
	OpTrackAlbionJournalAchievements: "TrackAlbionJournalAchievements",
	OpTrackAlbionJournalAchievementSubCategory: "TrackAlbionJournalAchievementSubCategory",
	OpRequestOutlandsTeleportationUsage: "RequestOutlandsTeleportationUsage",
	OpPickupFromPiledObjectStart: "PickupFromPiledObjectStart",
	OpPickupFromPiledObjectCancel: "PickupFromPiledObjectCancel",
	OpAssetOverview: "AssetOverview",
	OpAssetOverviewTabs: "AssetOverviewTabs",
	OpAssetOverviewTabContent: "AssetOverviewTabContent",
	OpAssetOverviewUnfreezeCache: "AssetOverviewUnfreezeCache",
	OpAssetOverviewSearch: "AssetOverviewSearch",
	OpAssetOverviewSearchTabs: "AssetOverviewSearchTabs",
	OpAssetOverviewSearchTabContent: "AssetOverviewSearchTabContent",
	OpAssetOverviewRecoverPlayerVault: "AssetOverviewRecoverPlayerVault",
	OpImmortalizeKillTrophy: "ImmortalizeKillTrophy",
	OpArmorySearch: "ArmorySearch",
	OpArmoryItemUsageStatistics: "ArmoryItemUsageStatistics",
	OpArmoryActivityUsageStatistics: "ArmoryActivityUsageStatistics",
	OpHellDungeonUseStaticEntrance: "HellDungeonUseStaticEntrance",
	OpTravelIslandShowroom: "TravelIslandShowroom",
	OpGetXuids: "GetXuids",
	OpXboxServiceTicket: "XboxServiceTicket",
	OpEvaluatePlatformPerks: "EvaluatePlatformPerks",
	OpLinkAccountToXbox: "LinkAccountToXbox",
	OpTravelFactionWarfarePortal: "TravelFactionWarfarePortal",
	OpRequestRedZoneEventStandings: "RequestRedZoneEventStandings",
	OpGetZergDebuffInfo: "GetZergDebuffInfo",
	OpRequestLoreSnippetStates: "RequestLoreSnippetStates",
	OpRetrieveCarriableObjectStart: "RetrieveCarriableObjectStart",
	OpRetrieveCarriableObjectCancel: "RetrieveCarriableObjectCancel",
}

func (c OperationCode) String() string { if n,ok:=operationCodeNames[c]; ok { return n }; return "Operation("+itoa(int(c))+")" }
func IsKnownOperationCode(c OperationCode) bool { _,ok:=operationCodeNames[c]; return ok }

