// Code generated. DO NOT EDIT. Albion Photon event codes (param 252).
package photon

// EventCode is the Photon event code (param 252).
type EventCode int16

const (
	EvLeave EventCode = 1
	EvJoinFinished EventCode = 2
	EvMove EventCode = 3
	EvTeleport EventCode = 4
	EvChangeEquipment EventCode = 5
	EvHealthUpdate EventCode = 6
	EvHealthUpdates EventCode = 7
	EvEnergyUpdate EventCode = 8
	EvDamageShieldUpdate EventCode = 9
	EvCraftingFocusUpdate EventCode = 10
	EvActiveSpellEffectsUpdate EventCode = 11
	EvResetCooldowns EventCode = 12
	EvAttack EventCode = 13
	EvCastStart EventCode = 14
	EvChannelingUpdate EventCode = 15
	EvCastCancel EventCode = 16
	EvCastTimeUpdate EventCode = 17
	EvCastFinished EventCode = 18
	EvCastSpell EventCode = 19
	EvCastSpells EventCode = 20
	EvCastHit EventCode = 21
	EvCastHits EventCode = 22
	EvStoredTargetsUpdate EventCode = 23
	EvChannelingEnded EventCode = 24
	EvAttackBuilding EventCode = 25
	EvInventoryPutItem EventCode = 26
	EvInventoryDeleteItem EventCode = 27
	EvInventoryState EventCode = 28
	EvNewCharacter EventCode = 29
	EvNewEquipmentItem EventCode = 30
	EvNewSiegeBannerItem EventCode = 31
	EvNewSimpleItem EventCode = 32
	EvNewFurnitureItem EventCode = 33
	EvNewKillTrophyItem EventCode = 34
	EvNewJournalItem EventCode = 35
	EvNewLaborerItem EventCode = 36
	EvNewEquipmentItemLegendarySoul EventCode = 37
	EvNewSimpleHarvestableObject EventCode = 38
	EvNewSimpleHarvestableObjectList EventCode = 39
	EvNewHarvestableObject EventCode = 40
	EvNewTreasureDestinationObject EventCode = 41
	EvTreasureDestinationObjectStatus EventCode = 42
	EvCloseTreasureDestinationObject EventCode = 43
	EvNewSilverObject EventCode = 44
	EvNewBuilding EventCode = 45
	EvHarvestableChangeState EventCode = 46
	EvMobChangeState EventCode = 47
	EvFactionBuildingInfo EventCode = 48
	EvCraftBuildingInfo EventCode = 49
	EvRepairBuildingInfo EventCode = 50
	EvMeldBuildingInfo EventCode = 51
	EvConstructionSiteInfo EventCode = 52
	EvPlayerBuildingInfo EventCode = 53
	EvFarmBuildingInfo EventCode = 54
	EvTutorialBuildingInfo EventCode = 55
	EvLaborerObjectInfo EventCode = 56
	EvLaborerObjectJobInfo EventCode = 57
	EvMarketPlaceBuildingInfo EventCode = 58
	EvHarvestStart EventCode = 59
	EvHarvestCancel EventCode = 60
	EvHarvestFinished EventCode = 61
	EvTakeSilver EventCode = 62
	EvRemoveSilver EventCode = 63
	EvActionOnBuildingStart EventCode = 64
	EvActionOnBuildingCancel EventCode = 65
	EvActionOnBuildingFinished EventCode = 66
	EvItemRerollQualityFinished EventCode = 67
	EvInstallResourceStart EventCode = 68
	EvInstallResourceCancel EventCode = 69
	EvInstallResourceFinished EventCode = 70
	EvCraftItemFinished EventCode = 71
	EvLogoutCancel EventCode = 72
	EvChatMessage EventCode = 73
	EvChatSay EventCode = 74
	EvChatWhisper EventCode = 75
	EvChatMuted EventCode = 76
	EvPlayEmote EventCode = 77
	EvStopEmote EventCode = 78
	EvSystemMessage EventCode = 79
	EvUtilityTextMessage EventCode = 80
	EvUpdateMoney EventCode = 81
	EvUpdateFame EventCode = 82
	EvUpdateLearningPoints EventCode = 83
	EvUpdateReSpecPoints EventCode = 84
	EvUpdateCurrency EventCode = 85
	EvUpdateFactionStanding EventCode = 86
	EvUpdateStanding EventCode = 87
	EvRespawn EventCode = 88
	EvServerDebugLog EventCode = 89
	EvCharacterEquipmentChanged EventCode = 90
	EvRegenerationHealthChanged EventCode = 91
	EvRegenerationEnergyChanged EventCode = 92
	EvRegenerationMountHealthChanged EventCode = 93
	EvRegenerationCraftingChanged EventCode = 94
	EvRegenerationHealthEnergyComboChanged EventCode = 95
	EvRegenerationPlayerComboChanged EventCode = 96
	EvDurabilityChanged EventCode = 97
	EvNewLoot EventCode = 98
	EvAttachItemContainer EventCode = 99
	EvDetachItemContainer EventCode = 100
	EvInvalidateItemContainer EventCode = 101
	EvLockItemContainer EventCode = 102
	EvGuildUpdate EventCode = 103
	EvGuildPlayerUpdated EventCode = 104
	EvInvitedToGuild EventCode = 105
	EvGuildMemberWorldUpdate EventCode = 106
	EvUpdateMatchDetails EventCode = 107
	EvObjectEvent EventCode = 108
	EvNewMonolithObject EventCode = 109
	EvMonolithHasBannersPlacedUpdate EventCode = 110
	EvNewOrbObject EventCode = 111
	EvNewCastleObject EventCode = 112
	EvNewSpellEffectArea EventCode = 113
	EvUpdateSpellEffectArea EventCode = 114
	EvNewChainSpell EventCode = 115
	EvUpdateChainSpell EventCode = 116
	EvNewTreasureChest EventCode = 117
	EvStartMatch EventCode = 118
	EvStartArenaMatchInfos EventCode = 119
	EvEndArenaMatch EventCode = 120
	EvMatchUpdate EventCode = 121
	EvActiveMatchUpdate EventCode = 122
	EvNewMob EventCode = 123
	EvDebugMobInfo EventCode = 124
	EvDebugVariablesInfo EventCode = 125
	EvDebugReputationInfo EventCode = 126
	EvDebugDiminishingReturnInfo EventCode = 127
	EvDebugSmartClusterQueueInfo EventCode = 128
	EvClaimOrbStart EventCode = 129
	EvClaimOrbFinished EventCode = 130
	EvClaimOrbCancel EventCode = 131
	EvOrbUpdate EventCode = 132
	EvOrbClaimed EventCode = 133
	EvOrbReset EventCode = 134
	EvNewWarCampObject EventCode = 135
	EvNewMatchLootChestObject EventCode = 136
	EvNewArenaExit EventCode = 137
	EvGuildMemberTerritoryUpdate EventCode = 138
	EvInvitedMercenaryToMatch EventCode = 139
	EvClusterInfoUpdate EventCode = 140
	EvForcedMovement EventCode = 141
	EvForcedMovementCancel EventCode = 142
	EvCharacterStats EventCode = 143
	EvCharacterStatsKillHistory EventCode = 144
	EvCharacterStatsDeathHistory EventCode = 145
	EvCharacterStatsKnockDownHistory EventCode = 146
	EvCharacterStatsKnockedDownHistory EventCode = 147
	EvGuildStats EventCode = 148
	EvKillHistoryDetails EventCode = 149
	EvItemKillHistoryDetails EventCode = 150
	EvFullAchievementInfo EventCode = 151
	EvFinishedAchievement EventCode = 152
	EvAchievementProgressInfo EventCode = 153
	EvFullAchievementProgressInfo EventCode = 154
	EvFullTrackedAchievementInfo EventCode = 155
	EvFullAutoLearnAchievementInfo EventCode = 156
	EvQuestGiverQuestOffered EventCode = 157
	EvQuestGiverDebugInfo EventCode = 158
	EvConsoleEvent EventCode = 159
	EvTimeSync EventCode = 160
	EvChangeAvatar EventCode = 161
	EvChangeMountSkin EventCode = 162
	EvGameEvent EventCode = 163
	EvKilledPlayer EventCode = 164
	EvDied EventCode = 165
	EvKnockedDown EventCode = 166
	EvUnconcious EventCode = 167
	EvMatchPlayerJoinedEvent EventCode = 168
	EvMatchPlayerStatsEvent EventCode = 169
	EvMatchPlayerStatsCompleteEvent EventCode = 170
	EvMatchTimeLineEventEvent EventCode = 171
	EvMatchNewCombatRound EventCode = 172
	EvMatchEndCombatRound EventCode = 173
	EvMatchPlayerMainGearStatsEvent EventCode = 174
	EvMatchPlayerChangedAvatarEvent EventCode = 175
	EvInvitationPlayerTrade EventCode = 176
	EvPlayerTradeStart EventCode = 177
	EvPlayerTradeCancel EventCode = 178
	EvPlayerTradeUpdate EventCode = 179
	EvPlayerTradeFinished EventCode = 180
	EvPlayerTradeAcceptChange EventCode = 181
	EvMiniMapPing EventCode = 182
	EvMarketPlaceNotification EventCode = 183
	EvDuellingChallengePlayer EventCode = 184
	EvNewDuellingPost EventCode = 185
	EvDuelStarted EventCode = 186
	EvDuelEnded EventCode = 187
	EvDuelDenied EventCode = 188
	EvDuelRequestCanceled EventCode = 189
	EvDuelLeftArea EventCode = 190
	EvDuelReEnteredArea EventCode = 191
	EvNewRealEstate EventCode = 192
	EvMiniMapOwnedBuildingsPositions EventCode = 193
	EvRealEstateListUpdate EventCode = 194
	EvGuildLogoUpdate EventCode = 195
	EvGuildLogoChanged EventCode = 196
	EvPlaceableObjectPlace EventCode = 197
	EvPlaceableObjectPlaceCancel EventCode = 198
	EvFurnitureObjectBuffProviderInfo EventCode = 199
	EvFurnitureObjectCheatProviderInfo EventCode = 200
	EvFarmableObjectInfo EventCode = 201
	EvNewUnreadMails EventCode = 202
	EvMailOperationPossible EventCode = 203
	EvGuildLogoObjectUpdate EventCode = 204
	EvStartLogout EventCode = 205
	EvNewChatChannels EventCode = 206
	EvJoinedChatChannel EventCode = 207
	EvLeftChatChannel EventCode = 208
	EvRemovedChatChannel EventCode = 209
	EvAccessStatus EventCode = 210
	EvMounted EventCode = 211
	EvMountStart EventCode = 212
	EvMountCancel EventCode = 213
	EvNewTravelpoint EventCode = 214
	EvNewIslandAccessPoint EventCode = 215
	EvNewExit EventCode = 216
	EvUpdateHome EventCode = 217
	EvUpdateChatSettings EventCode = 218
	EvResurrectionOffer EventCode = 219
	EvResurrectionReply EventCode = 220
	EvLootEquipmentChanged EventCode = 221
	EvUpdateUnlockedGuildLogos EventCode = 222
	EvUpdateUnlockedAvatars EventCode = 223
	EvUpdateUnlockedAvatarRings EventCode = 224
	EvUpdateUnlockedBuildings EventCode = 225
	EvNewIslandManagement EventCode = 226
	EvNewTeleportStone EventCode = 227
	EvCloak EventCode = 228
	EvPartyInvitation EventCode = 229
	EvPartyJoinRequest EventCode = 230
	EvPartyJoined EventCode = 231
	EvPartyDisbanded EventCode = 232
	EvPartyPlayerJoined EventCode = 233
	EvPartyChangedOrder EventCode = 234
	EvPartyPlayerLeft EventCode = 235
	EvPartyLeaderChanged EventCode = 236
	EvPartyLootSettingChangedPlayer EventCode = 237
	EvPartySilverGained EventCode = 238
	EvPartyPlayerUpdated EventCode = 239
	EvPartyInvitationAnswer EventCode = 240
	EvPartyJoinRequestAnswer EventCode = 241
	EvPartyMarkedObjectsUpdated EventCode = 242
	EvPartyOnClusterPartyJoined EventCode = 243
	EvPartySetRoleFlag EventCode = 244
	EvPartyInviteOrJoinPlayerEquipmentInfo EventCode = 245
	EvPartyReadyCheckUpdate EventCode = 246
	EvPartyFactionWarfareReinforcementSettingChangedPlayer EventCode = 247
	EvSpellCooldownUpdate EventCode = 248
	EvNewHellgateExitPortal EventCode = 249
	EvNewExpeditionExit EventCode = 250
	EvNewExpeditionNarrator EventCode = 251
	EvExitEnterStart EventCode = 252
	EvExitEnterCancel EventCode = 253
	EvExitEnterFinished EventCode = 254
	EvNewQuestGiverObject EventCode = 255
	EvFullQuestInfo EventCode = 256
	EvQuestProgressInfo EventCode = 257
	EvQuestGiverInfoForPlayer EventCode = 258
	EvFullExpeditionInfo EventCode = 259
	EvExpeditionQuestProgressInfo EventCode = 260
	EvInvitedToExpedition EventCode = 261
	EvExpeditionRegistrationInfo EventCode = 262
	EvEnteringExpeditionStart EventCode = 263
	EvEnteringExpeditionCancel EventCode = 264
	EvRewardGranted EventCode = 265
	EvArenaRegistrationInfo EventCode = 266
	EvEnteringArenaStart EventCode = 267
	EvEnteringArenaCancel EventCode = 268
	EvEnteringArenaLockStart EventCode = 269
	EvEnteringArenaLockCancel EventCode = 270
	EvInvitedToArenaMatch EventCode = 271
	EvUsingHellgateShrine EventCode = 272
	EvEnteringHellgateLockStart EventCode = 273
	EvEnteringHellgateLockCancel EventCode = 274
	EvPlayerCounts EventCode = 275
	EvInCombatStateUpdate EventCode = 276
	EvOtherGrabbedLoot EventCode = 277
	EvTreasureChestUsingStart EventCode = 278
	EvTreasureChestUsingFinished EventCode = 279
	EvTreasureChestUsingCancel EventCode = 280
	EvTreasureChestUsingOpeningComplete EventCode = 281
	EvTreasureChestForceCloseInventory EventCode = 282
	EvLocalTreasuresUpdate EventCode = 283
	EvLootChestSpawnpointsUpdate EventCode = 284
	EvPremiumChanged EventCode = 285
	EvPremiumExtended EventCode = 286
	EvPremiumLifeTimeRewardGained EventCode = 287
	EvGoldPurchased EventCode = 288
	EvLaborerGotUpgraded EventCode = 289
	EvJournalGotFull EventCode = 290
	EvJournalFillError EventCode = 291
	EvFriendRequest EventCode = 292
	EvFriendRequestInfos EventCode = 293
	EvFriendInfos EventCode = 294
	EvFriendRequestAnswered EventCode = 295
	EvFriendOnlineStatus EventCode = 296
	EvFriendRequestCanceled EventCode = 297
	EvFriendRemoved EventCode = 298
	EvFriendUpdated EventCode = 299
	EvPartyLootItems EventCode = 300
	EvPartyLootItemsRemoved EventCode = 301
	EvPartyLootItemTypesRemoved EventCode = 302
	EvReputationUpdate EventCode = 303
	EvDefenseUnitAttackBegin EventCode = 304
	EvDefenseUnitAttackEnd EventCode = 305
	EvDefenseUnitAttackDamage EventCode = 306
	EvUnrestrictedPvpZoneUpdate EventCode = 307
	EvUnrestrictedPvpZoneStatus EventCode = 308
	EvReputationImplicationUpdate EventCode = 309
	EvNewMountObject EventCode = 310
	EvMountHealthUpdate EventCode = 311
	EvMountCooldownUpdate EventCode = 312
	EvNewExpeditionAgent EventCode = 313
	EvNewExpeditionCheckPoint EventCode = 314
	EvExpeditionStartEvent EventCode = 315
	EvVoteEvent EventCode = 316
	EvRatingEvent EventCode = 317
	EvNewArenaAgent EventCode = 318
	EvBoostFarmable EventCode = 319
	EvUseFunction EventCode = 320
	EvNewPortalEntrance EventCode = 321
	EvNewPortalExit EventCode = 322
	EvNewRandomDungeonExit EventCode = 323
	EvWaitingQueueUpdate EventCode = 324
	EvPlayerMovementRateUpdate EventCode = 325
	EvObserveStart EventCode = 326
	EvMinimapZergs EventCode = 327
	EvMinimapSmartClusterZergs EventCode = 328
	EvPaymentTransactions EventCode = 329
	EvPerformanceStatsUpdate EventCode = 330
	EvOverloadModeUpdate EventCode = 331
	EvDebugDrawEvent EventCode = 332
	EvRecordCameraMove EventCode = 333
	EvRecordStart EventCode = 334
	EvDeliverCarriableObjectStart EventCode = 335
	EvDeliverCarriableObjectCancel EventCode = 336
	EvDeliverCarriableObjectReset EventCode = 337
	EvDeliverCarriableObjectFinished EventCode = 338
	EvTerritoryClaimStart EventCode = 339
	EvTerritoryClaimCancel EventCode = 340
	EvTerritoryClaimFinished EventCode = 341
	EvTerritoryScheduleResult EventCode = 342
	EvTerritoryUpgradeWithPowerCrystalResult EventCode = 343
	EvReceiveCarriableObjectStart EventCode = 344
	EvReceiveCarriableObjectFinished EventCode = 345
	EvUpdateAccountState EventCode = 346
	EvStartDeterministicRoam EventCode = 347
	EvGuildFullAccessTagsUpdated EventCode = 348
	EvGuildAccessTagUpdated EventCode = 349
	EvGvgSeasonUpdate EventCode = 350
	EvGvgSeasonCheatCommand EventCode = 351
	EvSeasonPointsByKillingBooster EventCode = 352
	EvFishingStart EventCode = 353
	EvFishingCast EventCode = 354
	EvFishingCatch EventCode = 355
	EvFishingFinished EventCode = 356
	EvFishingCancel EventCode = 357
	EvNewFloatObject EventCode = 358
	EvNewFishingZoneObject EventCode = 359
	EvFishingMiniGame EventCode = 360
	EvAlbionJournalAchievementCompleted EventCode = 361
	EvUpdatePuppet EventCode = 362
	EvChangeFlaggingFinished EventCode = 363
	EvNewOutpostObject EventCode = 364
	EvOutpostUpdate EventCode = 365
	EvOutpostClaimed EventCode = 366
	EvOverChargeEnd EventCode = 367
	EvOverChargeStatus EventCode = 368
	EvPartyFinderFullUpdate EventCode = 369
	EvPartyFinderUpdate EventCode = 370
	EvPartyFinderApplicantsUpdate EventCode = 371
	EvPartyFinderEquipmentSnapshot EventCode = 372
	EvPartyFinderJoinRequestDeclined EventCode = 373
	EvNewUnlockedPersonalSeasonRewards EventCode = 374
	EvPersonalSeasonPointsGained EventCode = 375
	EvPersonalSeasonPastSeasonDataEvent EventCode = 376
	EvMatchLootChestOpeningStart EventCode = 377
	EvMatchLootChestOpeningFinished EventCode = 378
	EvMatchLootChestOpeningCancel EventCode = 379
	EvNotifyCrystalMatchReward EventCode = 380
	EvCrystalRealmFeedback EventCode = 381
	EvNewLocationMarker EventCode = 382
	EvNewTutorialBlocker EventCode = 383
	EvNewTileSwitch EventCode = 384
	EvNewInformationProvider EventCode = 385
	EvNewDynamicGuildLogo EventCode = 386
	EvNewDecoration EventCode = 387
	EvTutorialUpdate EventCode = 388
	EvTriggerHintBox EventCode = 389
	EvRandomDungeonPositionInfo EventCode = 390
	EvNewLootChest EventCode = 391
	EvUpdateLootChest EventCode = 392
	EvLootChestOpened EventCode = 393
	EvUpdateLootProtectedByMobsWithMinimapDisplay EventCode = 394
	EvNewShrine EventCode = 395
	EvUpdateShrine EventCode = 396
	EvUpdateRoom EventCode = 397
	EvNewMobSoul EventCode = 398
	EvNewHellgateShrine EventCode = 399
	EvUpdateHellgateShrine EventCode = 400
	EvActivateHellgateExit EventCode = 401
	EvMutePlayerUpdate EventCode = 402
	EvShopTileUpdate EventCode = 403
	EvShopUpdate EventCode = 404
	EvAntiCheatKick EventCode = 405
	EvBattlEyeServerMessage EventCode = 406
	EvUnlockVanityUnlock EventCode = 407
	EvAvatarUnlocked EventCode = 408
	EvCustomizationChanged EventCode = 409
	EvBaseVaultInfo EventCode = 410
	EvGuildVaultInfo EventCode = 411
	EvBankVaultInfo EventCode = 412
	EvRecoveryVaultPlayerInfo EventCode = 413
	EvRecoveryVaultGuildInfo EventCode = 414
	EvUpdateWardrobe EventCode = 415
	EvCastlePhaseChanged EventCode = 416
	EvGuildAccountLogEvent EventCode = 417
	EvNewHideoutObject EventCode = 418
	EvNewHideoutManagement EventCode = 419
	EvNewHideoutExit EventCode = 420
	EvInitHideoutAttackStart EventCode = 421
	EvInitHideoutAttackCancel EventCode = 422
	EvInitHideoutAttackFinished EventCode = 423
	EvHideoutManagementUpdate EventCode = 424
	EvHideoutUpgradeWithPowerCrystalResult EventCode = 425
	EvIpChanged EventCode = 426
	EvSmartClusterQueueUpdateInfo EventCode = 427
	EvSmartClusterQueueActiveInfo EventCode = 428
	EvSmartClusterQueueKickWarning EventCode = 429
	EvSmartClusterQueueInvite EventCode = 430
	EvReceivedGvgSeasonPoints EventCode = 431
	EvTowerPowerPointUpdate EventCode = 432
	EvOpenWorldAttackScheduleStart EventCode = 433
	EvOpenWorldAttackScheduleFinished EventCode = 434
	EvOpenWorldAttackScheduleCancel EventCode = 435
	EvOpenWorldAttackConquerStart EventCode = 436
	EvOpenWorldAttackConquerFinished EventCode = 437
	EvOpenWorldAttackConquerCancel EventCode = 438
	EvOpenWorldAttackConquerStatus EventCode = 439
	EvOpenWorldAttackStart EventCode = 440
	EvOpenWorldAttackEnd EventCode = 441
	EvNewRandomResourceBlocker EventCode = 442
	EvNewHomeObject EventCode = 443
	EvHideoutObjectUpdate EventCode = 444
	EvUpdateInfamy EventCode = 445
	EvMinimapPositionMarkers EventCode = 446
	EvNewTunnelExit EventCode = 447
	EvCorruptedDungeonUpdate EventCode = 448
	EvCorruptedDungeonStatus EventCode = 449
	EvCorruptedDungeonInfamy EventCode = 450
	EvHellgateRestrictedAreaUpdate EventCode = 451
	EvHellgateInfamy EventCode = 452
	EvHellgateStatus EventCode = 453
	EvHellgateStatusUpdate EventCode = 454
	EvHellgateSuspense EventCode = 455
	EvReplaceSpellSlotWithMultiSpell EventCode = 456
	EvNewCorruptedShrine EventCode = 457
	EvUpdateCorruptedShrine EventCode = 458
	EvCorruptedShrineUsageStart EventCode = 459
	EvCorruptedShrineUsageCancel EventCode = 460
	EvExitUsed EventCode = 461
	EvLinkedToObject EventCode = 462
	EvLinkToObjectBroken EventCode = 463
	EvEstimatedMarketValueUpdate EventCode = 464
	EvStuckCancel EventCode = 465
	EvDungonEscapeReady EventCode = 466
	EvFactionWarfareClusterState EventCode = 467
	EvFactionWarfareHasUnclaimedWeeklyReportsEvent EventCode = 468
	EvSimpleFeedback EventCode = 469
	EvSmartClusterQueueSkipClusterError EventCode = 470
	EvXignCodeEvent EventCode = 471
	EvBatchUseItemStart EventCode = 472
	EvBatchUseItemEnd EventCode = 473
	EvRedZonePlayerNotification EventCode = 474
	EvRedZoneEventCheatCleanup EventCode = 475
	EvRedZoneFortressEventChestOpened EventCode = 476
	EvRedZoneWorldMapEvent EventCode = 477
	EvFactionWarfareStats EventCode = 478
	EvUpdateFactionBalanceFactors EventCode = 479
	EvFactionEnlistmentChanged EventCode = 480
	EvUpdateFactionRank EventCode = 481
	EvFactionWarfareCampaignRewardsUnlocked EventCode = 482
	EvFeaturedFeatureUpdate EventCode = 483
	EvNewCarriableObject EventCode = 484
	EvMinimapCrystalPositionMarker EventCode = 485
	EvCarriedObjectUpdate EventCode = 486
	EvPickupCarriableObjectStart EventCode = 487
	EvPickupCarriableObjectCancel EventCode = 488
	EvPickupCarriableObjectFinished EventCode = 489
	EvDoSimpleActionStart EventCode = 490
	EvDoSimpleActionCancel EventCode = 491
	EvDoSimpleActionFinished EventCode = 492
	EvNotifyGuestAccountVerified EventCode = 493
	EvMightAndFavorReceivedEvent EventCode = 494
	EvWeeklyPvpChallengeRewardStateUpdate EventCode = 495
	EvNewUnlockedPvpSeasonChallengeRewards EventCode = 496
	EvStaticDungeonEntrancesDungeonEventStatusUpdates EventCode = 497
	EvStaticDungeonDungeonValueUpdate EventCode = 498
	EvStaticDungeonEntranceDungeonEventsAborted EventCode = 499
	EvInAppPurchaseConfirmedGooglePlay EventCode = 500
	EvFeatureSwitchInfo EventCode = 501
	EvPartyJoinRequestAborted EventCode = 502
	EvPartyInviteAborted EventCode = 503
	EvPartyStartHuntRequest EventCode = 504
	EvPartyStartHuntRequested EventCode = 505
	EvPartyStartHuntRequestAnswer EventCode = 506
	EvPartyPlayerLeaveScheduled EventCode = 507
	EvGuildInviteDeclined EventCode = 508
	EvCancelMultiSpellSlots EventCode = 509
	EvNewVisualEventObject EventCode = 510
	EvCastleClaimProgress EventCode = 511
	EvCastleClaimProgressLogo EventCode = 512
	EvTownPortalUpdateState EventCode = 513
	EvTownPortalFailed EventCode = 514
	EvConsumableVanityChargesAdded EventCode = 515
	EvFestivitiesUpdate EventCode = 516
	EvNewBannerObject EventCode = 517
	EvNewMistsImmediateReturnExit EventCode = 518
	EvMistsPlayerJoinedInfo EventCode = 519
	EvNewMistsStaticEntrance EventCode = 520
	EvNewMistsOpenWorldExit EventCode = 521
	EvNewTunnelExitTemp EventCode = 522
	EvNewMistsWispSpawn EventCode = 523
	EvMistsWispSpawnStateChange EventCode = 524
	EvNewMistsCityEntrance EventCode = 525
	EvNewMistsCityRoadsEntrance EventCode = 526
	EvMistsCityRoadsEntrancePartyStateUpdate EventCode = 527
	EvMistsCityRoadsEntranceClearStateForParty EventCode = 528
	EvMistsEntranceDataChanged EventCode = 529
	EvNewCagedObject EventCode = 530
	EvCagedObjectStateUpdated EventCode = 531
	EvEntrancePartyBindingCreated EventCode = 532
	EvEntrancePartyBindingCleared EventCode = 533
	EvEntrancePartyBindingInfos EventCode = 534
	EvNewMistsBorderExit EventCode = 535
	EvNewMistsDungeonExit EventCode = 536
	EvLocalQuestInfos EventCode = 537
	EvLocalQuestStarted EventCode = 538
	EvLocalQuestActive EventCode = 539
	EvLocalQuestInactive EventCode = 540
	EvLocalQuestProgressUpdate EventCode = 541
	EvNewUnrestrictedPvpZone EventCode = 542
	EvTemporaryFlaggingStatusUpdate EventCode = 543
	EvSpellTestPerformanceUpdate EventCode = 544
	EvTransformation EventCode = 545
	EvTransformationEnd EventCode = 546
	EvUpdateTrustlevel EventCode = 547
	EvRevealHiddenTimeStamps EventCode = 548
	EvModifyItemTraitFinished EventCode = 549
	EvRerollItemTraitValueFinished EventCode = 550
	EvHuntQuestProgressInfo EventCode = 551
	EvHuntStarted EventCode = 552
	EvHuntFinished EventCode = 553
	EvHuntAborted EventCode = 554
	EvHuntMissionStepStateUpdate EventCode = 555
	EvNewHuntTrack EventCode = 556
	EvHuntMissionUpdate EventCode = 557
	EvHuntQuestMissionProgressUpdate EventCode = 558
	EvHuntTrackUsed EventCode = 559
	EvHuntTrackUseableAgain EventCode = 560
	EvMinimapHuntTrackMarkers EventCode = 561
	EvNoTracksFound EventCode = 562
	EvHuntQuestAborted EventCode = 563
	EvInteractWithTrackStart EventCode = 564
	EvInteractWithTrackCancel EventCode = 565
	EvInteractWithTrackFinished EventCode = 566
	EvNewDynamicCompound EventCode = 567
	EvLegendaryItemDestroyed EventCode = 568
	EvAttunementInfo EventCode = 569
	EvTerritoryClaimRaidedRawEnergyCrystalResult EventCode = 570
	EvCarriedObjectExpiryWarning EventCode = 571
	EvCarriedObjectExpired EventCode = 572
	EvTerritoryRaidStart EventCode = 573
	EvTerritoryRaidCancel EventCode = 574
	EvTerritoryRaidFinished EventCode = 575
	EvTerritoryRaidResult EventCode = 576
	EvTerritoryMonolithActiveRaidStatus EventCode = 577
	EvTerritoryMonolithActiveRaidCancelled EventCode = 578
	EvMonolithEnergyStorageUpdate EventCode = 579
	EvMonolithNextScheduledOpenWorldAttackUpdate EventCode = 580
	EvMonolithProtectedBuildingsDamageReductionUpdate EventCode = 581
	EvNewBuildingBaseEvent EventCode = 582
	EvNewFortificationBuilding EventCode = 583
	EvNewCastleGateBuilding EventCode = 584
	EvBuildingDurabilityUpdate EventCode = 585
	EvMonolithFortificationPointsUpdate EventCode = 586
	EvFortificationBuildingUpgradeInfo EventCode = 587
	EvFortificationBuildingsDamageStateUpdate EventCode = 588
	EvSiegeNotificationEvent EventCode = 589
	EvUpdateEnemyWarBannerActive EventCode = 590
	EvTerritoryAnnouncePlayerEjection EventCode = 591
	EvCastleGateSwitchUseStarted EventCode = 592
	EvCastleGateSwitchUseFinished EventCode = 593
	EvFortificationBuildingWillDowngrade EventCode = 594
	EvBotCommand EventCode = 595
	EvJournalAchievementProgressUpdate EventCode = 596
	EvJournalClaimableRewardUpdate EventCode = 597
	EvKeySync EventCode = 598
	EvLocalQuestAreaGone EventCode = 599
	EvDynamicTemplate EventCode = 600
	EvDynamicTemplateForcedStateChange EventCode = 601
	EvNewOutlandsTeleportationPortal EventCode = 602
	EvNewOutlandsTeleportationReturnPortal EventCode = 603
	EvOutlandsTeleportationBindingCleared EventCode = 604
	EvOutlandsTeleportationReturnPortalUpdateEvent EventCode = 605
	EvPlayerUsedOutlandsTeleportationPortal EventCode = 606
	EvEncumberedRestricted EventCode = 607
	EvNewPiledObject EventCode = 608
	EvPiledObjectStateChanged EventCode = 609
	EvNewSmugglerCrateDeliveryStation EventCode = 610
	EvKillRewardedNoFame EventCode = 611
	EvPickupFromPiledObjectStart EventCode = 612
	EvPickupFromPiledObjectCancel EventCode = 613
	EvPickupFromPiledObjectReset EventCode = 614
	EvPickupFromPiledObjectFinished EventCode = 615
	EvArmoryActivityChange EventCode = 616
	EvNewKillTrophyFurnitureBuilding EventCode = 617
	EvHellDungeonsPlayerJoinedInfo EventCode = 618
	EvNewTileSwitchTrigger EventCode = 619
	EvNewMultiRewardObject EventCode = 620
	EvNewHellDungeonSoulShrineObject EventCode = 621
	EvHellDungeonSoulShrineStateUpdate EventCode = 622
	EvNewResurrectionShrine EventCode = 623
	EvUpdateResurrectionShrine EventCode = 624
	EvStandTimeFinished EventCode = 625
	EvEpicAchievementAndStatsUpdate EventCode = 626
	EvSpectateTargetAfterDeathUpdate EventCode = 627
	EvSpectateTargetAfterDeathEnded EventCode = 628
	EvNewHellDungeonUpwardExit EventCode = 629
	EvNewHellDungeonSoulExit EventCode = 630
	EvNewHellDungeonDownwardExit EventCode = 631
	EvNewHellDungeonChestExit EventCode = 632
	EvNewCorruptedStaticEntrance EventCode = 633
	EvNewHellDungeonStaticEntrance EventCode = 634
	EvUpdateHellDungeonStaticEntranceState EventCode = 635
	EvDebugTriggerHellDungeonShutdownStart EventCode = 636
	EvFullJournalQuestInfo EventCode = 637
	EvJournalQuestProgressInfo EventCode = 638
	EvNewHellDungeonRoomShrineObject EventCode = 639
	EvHellDungeonRoomShrineStateUpdate EventCode = 640
	EvSimpleBehaviourBuildingStateUpdate EventCode = 641
	EvSetTimeScaling EventCode = 642
	EvStopTimeScaling EventCode = 643
	EvKeyValidation EventCode = 644
	EvPlayerJoinMapMarkerTimerStates EventCode = 645
	EvNewMapMarkerTimer EventCode = 646
	EvRemoveMapMarkerTimer EventCode = 647
	EvNewFactionFortressObject EventCode = 648
	EvFactionFortressAnnouncePlayerEjection EventCode = 649
	EvRewardFactionWarfareSupply EventCode = 650
	EvFactionCaptureAreaProgressUpdate EventCode = 651
	EvFactionFortressClaimed EventCode = 652
	EvFactionFortressWeaponCachesSpawned EventCode = 653
	EvFactionFortressWeaponCacheClaimed EventCode = 654
	EvFactionFortressFightStateUpdate EventCode = 655
	EvFactionFortressCutoffFightStateUpdate EventCode = 656
	EvFactionFortressFightEnded EventCode = 657
	EvNewFactionWarfarePortal EventCode = 658
	EvFactionPortalTargetUpdate EventCode = 659
	EvFactionFortressFightStartedInRemoteClusterEvent EventCode = 660
	EvFactionFortressFightFinishedInRemoteClusterEvent EventCode = 661
	EvFactionDuchySupplyWarDefensiveVictoryEvent EventCode = 662
	EvFactionDuchyReconnectedFromCutoffEvent EventCode = 663
	EvFactionFortressCutoffFightCancelledByClusterOwnerChangeEvent EventCode = 664
	EvFactionDuchyEnteredCutoffStateEvent EventCode = 665
	EvLeaveProtectionStateUpdate EventCode = 666
	EvRedZoneEventStandings EventCode = 667
	EvNewFactionBattleStandardDeliveryStation EventCode = 668
	EvNewLoreSnippetObject EventCode = 669
	EvLoreSnippetObjectStateUpdate EventCode = 670
	EvLoreSnippedClaimed EventCode = 671
	EvLoreSnippetStatesChangedByCheat EventCode = 672
	EvNewTeleporterNode EventCode = 673
	EvTeleporterNodeStateChanged EventCode = 674
	EvTeleporterConnectionsFullStateUpdate EventCode = 675
	EvTeleporterConnectionStateChanged EventCode = 676
	EvRetrieveCarriableObjectStart EventCode = 677
	EvRetrieveCarriableObjectCancel EventCode = 678
	EvRetrieveCarriableObjectReset EventCode = 679
	EvRetrieveCarriableObjectFinished EventCode = 680
	EvLosingCarriableObjectStart EventCode = 681
	EvLosingCarriableObjectFinished EventCode = 682
)

var eventCodeNames = map[EventCode]string{
	EvLeave: "Leave",
	EvJoinFinished: "JoinFinished",
	EvMove: "Move",
	EvTeleport: "Teleport",
	EvChangeEquipment: "ChangeEquipment",
	EvHealthUpdate: "HealthUpdate",
	EvHealthUpdates: "HealthUpdates",
	EvEnergyUpdate: "EnergyUpdate",
	EvDamageShieldUpdate: "DamageShieldUpdate",
	EvCraftingFocusUpdate: "CraftingFocusUpdate",
	EvActiveSpellEffectsUpdate: "ActiveSpellEffectsUpdate",
	EvResetCooldowns: "ResetCooldowns",
	EvAttack: "Attack",
	EvCastStart: "CastStart",
	EvChannelingUpdate: "ChannelingUpdate",
	EvCastCancel: "CastCancel",
	EvCastTimeUpdate: "CastTimeUpdate",
	EvCastFinished: "CastFinished",
	EvCastSpell: "CastSpell",
	EvCastSpells: "CastSpells",
	EvCastHit: "CastHit",
	EvCastHits: "CastHits",
	EvStoredTargetsUpdate: "StoredTargetsUpdate",
	EvChannelingEnded: "ChannelingEnded",
	EvAttackBuilding: "AttackBuilding",
	EvInventoryPutItem: "InventoryPutItem",
	EvInventoryDeleteItem: "InventoryDeleteItem",
	EvInventoryState: "InventoryState",
	EvNewCharacter: "NewCharacter",
	EvNewEquipmentItem: "NewEquipmentItem",
	EvNewSiegeBannerItem: "NewSiegeBannerItem",
	EvNewSimpleItem: "NewSimpleItem",
	EvNewFurnitureItem: "NewFurnitureItem",
	EvNewKillTrophyItem: "NewKillTrophyItem",
	EvNewJournalItem: "NewJournalItem",
	EvNewLaborerItem: "NewLaborerItem",
	EvNewEquipmentItemLegendarySoul: "NewEquipmentItemLegendarySoul",
	EvNewSimpleHarvestableObject: "NewSimpleHarvestableObject",
	EvNewSimpleHarvestableObjectList: "NewSimpleHarvestableObjectList",
	EvNewHarvestableObject: "NewHarvestableObject",
	EvNewTreasureDestinationObject: "NewTreasureDestinationObject",
	EvTreasureDestinationObjectStatus: "TreasureDestinationObjectStatus",
	EvCloseTreasureDestinationObject: "CloseTreasureDestinationObject",
	EvNewSilverObject: "NewSilverObject",
	EvNewBuilding: "NewBuilding",
	EvHarvestableChangeState: "HarvestableChangeState",
	EvMobChangeState: "MobChangeState",
	EvFactionBuildingInfo: "FactionBuildingInfo",
	EvCraftBuildingInfo: "CraftBuildingInfo",
	EvRepairBuildingInfo: "RepairBuildingInfo",
	EvMeldBuildingInfo: "MeldBuildingInfo",
	EvConstructionSiteInfo: "ConstructionSiteInfo",
	EvPlayerBuildingInfo: "PlayerBuildingInfo",
	EvFarmBuildingInfo: "FarmBuildingInfo",
	EvTutorialBuildingInfo: "TutorialBuildingInfo",
	EvLaborerObjectInfo: "LaborerObjectInfo",
	EvLaborerObjectJobInfo: "LaborerObjectJobInfo",
	EvMarketPlaceBuildingInfo: "MarketPlaceBuildingInfo",
	EvHarvestStart: "HarvestStart",
	EvHarvestCancel: "HarvestCancel",
	EvHarvestFinished: "HarvestFinished",
	EvTakeSilver: "TakeSilver",
	EvRemoveSilver: "RemoveSilver",
	EvActionOnBuildingStart: "ActionOnBuildingStart",
	EvActionOnBuildingCancel: "ActionOnBuildingCancel",
	EvActionOnBuildingFinished: "ActionOnBuildingFinished",
	EvItemRerollQualityFinished: "ItemRerollQualityFinished",
	EvInstallResourceStart: "InstallResourceStart",
	EvInstallResourceCancel: "InstallResourceCancel",
	EvInstallResourceFinished: "InstallResourceFinished",
	EvCraftItemFinished: "CraftItemFinished",
	EvLogoutCancel: "LogoutCancel",
	EvChatMessage: "ChatMessage",
	EvChatSay: "ChatSay",
	EvChatWhisper: "ChatWhisper",
	EvChatMuted: "ChatMuted",
	EvPlayEmote: "PlayEmote",
	EvStopEmote: "StopEmote",
	EvSystemMessage: "SystemMessage",
	EvUtilityTextMessage: "UtilityTextMessage",
	EvUpdateMoney: "UpdateMoney",
	EvUpdateFame: "UpdateFame",
	EvUpdateLearningPoints: "UpdateLearningPoints",
	EvUpdateReSpecPoints: "UpdateReSpecPoints",
	EvUpdateCurrency: "UpdateCurrency",
	EvUpdateFactionStanding: "UpdateFactionStanding",
	EvUpdateStanding: "UpdateStanding",
	EvRespawn: "Respawn",
	EvServerDebugLog: "ServerDebugLog",
	EvCharacterEquipmentChanged: "CharacterEquipmentChanged",
	EvRegenerationHealthChanged: "RegenerationHealthChanged",
	EvRegenerationEnergyChanged: "RegenerationEnergyChanged",
	EvRegenerationMountHealthChanged: "RegenerationMountHealthChanged",
	EvRegenerationCraftingChanged: "RegenerationCraftingChanged",
	EvRegenerationHealthEnergyComboChanged: "RegenerationHealthEnergyComboChanged",
	EvRegenerationPlayerComboChanged: "RegenerationPlayerComboChanged",
	EvDurabilityChanged: "DurabilityChanged",
	EvNewLoot: "NewLoot",
	EvAttachItemContainer: "AttachItemContainer",
	EvDetachItemContainer: "DetachItemContainer",
	EvInvalidateItemContainer: "InvalidateItemContainer",
	EvLockItemContainer: "LockItemContainer",
	EvGuildUpdate: "GuildUpdate",
	EvGuildPlayerUpdated: "GuildPlayerUpdated",
	EvInvitedToGuild: "InvitedToGuild",
	EvGuildMemberWorldUpdate: "GuildMemberWorldUpdate",
	EvUpdateMatchDetails: "UpdateMatchDetails",
	EvObjectEvent: "ObjectEvent",
	EvNewMonolithObject: "NewMonolithObject",
	EvMonolithHasBannersPlacedUpdate: "MonolithHasBannersPlacedUpdate",
	EvNewOrbObject: "NewOrbObject",
	EvNewCastleObject: "NewCastleObject",
	EvNewSpellEffectArea: "NewSpellEffectArea",
	EvUpdateSpellEffectArea: "UpdateSpellEffectArea",
	EvNewChainSpell: "NewChainSpell",
	EvUpdateChainSpell: "UpdateChainSpell",
	EvNewTreasureChest: "NewTreasureChest",
	EvStartMatch: "StartMatch",
	EvStartArenaMatchInfos: "StartArenaMatchInfos",
	EvEndArenaMatch: "EndArenaMatch",
	EvMatchUpdate: "MatchUpdate",
	EvActiveMatchUpdate: "ActiveMatchUpdate",
	EvNewMob: "NewMob",
	EvDebugMobInfo: "DebugMobInfo",
	EvDebugVariablesInfo: "DebugVariablesInfo",
	EvDebugReputationInfo: "DebugReputationInfo",
	EvDebugDiminishingReturnInfo: "DebugDiminishingReturnInfo",
	EvDebugSmartClusterQueueInfo: "DebugSmartClusterQueueInfo",
	EvClaimOrbStart: "ClaimOrbStart",
	EvClaimOrbFinished: "ClaimOrbFinished",
	EvClaimOrbCancel: "ClaimOrbCancel",
	EvOrbUpdate: "OrbUpdate",
	EvOrbClaimed: "OrbClaimed",
	EvOrbReset: "OrbReset",
	EvNewWarCampObject: "NewWarCampObject",
	EvNewMatchLootChestObject: "NewMatchLootChestObject",
	EvNewArenaExit: "NewArenaExit",
	EvGuildMemberTerritoryUpdate: "GuildMemberTerritoryUpdate",
	EvInvitedMercenaryToMatch: "InvitedMercenaryToMatch",
	EvClusterInfoUpdate: "ClusterInfoUpdate",
	EvForcedMovement: "ForcedMovement",
	EvForcedMovementCancel: "ForcedMovementCancel",
	EvCharacterStats: "CharacterStats",
	EvCharacterStatsKillHistory: "CharacterStatsKillHistory",
	EvCharacterStatsDeathHistory: "CharacterStatsDeathHistory",
	EvCharacterStatsKnockDownHistory: "CharacterStatsKnockDownHistory",
	EvCharacterStatsKnockedDownHistory: "CharacterStatsKnockedDownHistory",
	EvGuildStats: "GuildStats",
	EvKillHistoryDetails: "KillHistoryDetails",
	EvItemKillHistoryDetails: "ItemKillHistoryDetails",
	EvFullAchievementInfo: "FullAchievementInfo",
	EvFinishedAchievement: "FinishedAchievement",
	EvAchievementProgressInfo: "AchievementProgressInfo",
	EvFullAchievementProgressInfo: "FullAchievementProgressInfo",
	EvFullTrackedAchievementInfo: "FullTrackedAchievementInfo",
	EvFullAutoLearnAchievementInfo: "FullAutoLearnAchievementInfo",
	EvQuestGiverQuestOffered: "QuestGiverQuestOffered",
	EvQuestGiverDebugInfo: "QuestGiverDebugInfo",
	EvConsoleEvent: "ConsoleEvent",
	EvTimeSync: "TimeSync",
	EvChangeAvatar: "ChangeAvatar",
	EvChangeMountSkin: "ChangeMountSkin",
	EvGameEvent: "GameEvent",
	EvKilledPlayer: "KilledPlayer",
	EvDied: "Died",
	EvKnockedDown: "KnockedDown",
	EvUnconcious: "Unconcious",
	EvMatchPlayerJoinedEvent: "MatchPlayerJoinedEvent",
	EvMatchPlayerStatsEvent: "MatchPlayerStatsEvent",
	EvMatchPlayerStatsCompleteEvent: "MatchPlayerStatsCompleteEvent",
	EvMatchTimeLineEventEvent: "MatchTimeLineEventEvent",
	EvMatchNewCombatRound: "MatchNewCombatRound",
	EvMatchEndCombatRound: "MatchEndCombatRound",
	EvMatchPlayerMainGearStatsEvent: "MatchPlayerMainGearStatsEvent",
	EvMatchPlayerChangedAvatarEvent: "MatchPlayerChangedAvatarEvent",
	EvInvitationPlayerTrade: "InvitationPlayerTrade",
	EvPlayerTradeStart: "PlayerTradeStart",
	EvPlayerTradeCancel: "PlayerTradeCancel",
	EvPlayerTradeUpdate: "PlayerTradeUpdate",
	EvPlayerTradeFinished: "PlayerTradeFinished",
	EvPlayerTradeAcceptChange: "PlayerTradeAcceptChange",
	EvMiniMapPing: "MiniMapPing",
	EvMarketPlaceNotification: "MarketPlaceNotification",
	EvDuellingChallengePlayer: "DuellingChallengePlayer",
	EvNewDuellingPost: "NewDuellingPost",
	EvDuelStarted: "DuelStarted",
	EvDuelEnded: "DuelEnded",
	EvDuelDenied: "DuelDenied",
	EvDuelRequestCanceled: "DuelRequestCanceled",
	EvDuelLeftArea: "DuelLeftArea",
	EvDuelReEnteredArea: "DuelReEnteredArea",
	EvNewRealEstate: "NewRealEstate",
	EvMiniMapOwnedBuildingsPositions: "MiniMapOwnedBuildingsPositions",
	EvRealEstateListUpdate: "RealEstateListUpdate",
	EvGuildLogoUpdate: "GuildLogoUpdate",
	EvGuildLogoChanged: "GuildLogoChanged",
	EvPlaceableObjectPlace: "PlaceableObjectPlace",
	EvPlaceableObjectPlaceCancel: "PlaceableObjectPlaceCancel",
	EvFurnitureObjectBuffProviderInfo: "FurnitureObjectBuffProviderInfo",
	EvFurnitureObjectCheatProviderInfo: "FurnitureObjectCheatProviderInfo",
	EvFarmableObjectInfo: "FarmableObjectInfo",
	EvNewUnreadMails: "NewUnreadMails",
	EvMailOperationPossible: "MailOperationPossible",
	EvGuildLogoObjectUpdate: "GuildLogoObjectUpdate",
	EvStartLogout: "StartLogout",
	EvNewChatChannels: "NewChatChannels",
	EvJoinedChatChannel: "JoinedChatChannel",
	EvLeftChatChannel: "LeftChatChannel",
	EvRemovedChatChannel: "RemovedChatChannel",
	EvAccessStatus: "AccessStatus",
	EvMounted: "Mounted",
	EvMountStart: "MountStart",
	EvMountCancel: "MountCancel",
	EvNewTravelpoint: "NewTravelpoint",
	EvNewIslandAccessPoint: "NewIslandAccessPoint",
	EvNewExit: "NewExit",
	EvUpdateHome: "UpdateHome",
	EvUpdateChatSettings: "UpdateChatSettings",
	EvResurrectionOffer: "ResurrectionOffer",
	EvResurrectionReply: "ResurrectionReply",
	EvLootEquipmentChanged: "LootEquipmentChanged",
	EvUpdateUnlockedGuildLogos: "UpdateUnlockedGuildLogos",
	EvUpdateUnlockedAvatars: "UpdateUnlockedAvatars",
	EvUpdateUnlockedAvatarRings: "UpdateUnlockedAvatarRings",
	EvUpdateUnlockedBuildings: "UpdateUnlockedBuildings",
	EvNewIslandManagement: "NewIslandManagement",
	EvNewTeleportStone: "NewTeleportStone",
	EvCloak: "Cloak",
	EvPartyInvitation: "PartyInvitation",
	EvPartyJoinRequest: "PartyJoinRequest",
	EvPartyJoined: "PartyJoined",
	EvPartyDisbanded: "PartyDisbanded",
	EvPartyPlayerJoined: "PartyPlayerJoined",
	EvPartyChangedOrder: "PartyChangedOrder",
	EvPartyPlayerLeft: "PartyPlayerLeft",
	EvPartyLeaderChanged: "PartyLeaderChanged",
	EvPartyLootSettingChangedPlayer: "PartyLootSettingChangedPlayer",
	EvPartySilverGained: "PartySilverGained",
	EvPartyPlayerUpdated: "PartyPlayerUpdated",
	EvPartyInvitationAnswer: "PartyInvitationAnswer",
	EvPartyJoinRequestAnswer: "PartyJoinRequestAnswer",
	EvPartyMarkedObjectsUpdated: "PartyMarkedObjectsUpdated",
	EvPartyOnClusterPartyJoined: "PartyOnClusterPartyJoined",
	EvPartySetRoleFlag: "PartySetRoleFlag",
	EvPartyInviteOrJoinPlayerEquipmentInfo: "PartyInviteOrJoinPlayerEquipmentInfo",
	EvPartyReadyCheckUpdate: "PartyReadyCheckUpdate",
	EvPartyFactionWarfareReinforcementSettingChangedPlayer: "PartyFactionWarfareReinforcementSettingChangedPlayer",
	EvSpellCooldownUpdate: "SpellCooldownUpdate",
	EvNewHellgateExitPortal: "NewHellgateExitPortal",
	EvNewExpeditionExit: "NewExpeditionExit",
	EvNewExpeditionNarrator: "NewExpeditionNarrator",
	EvExitEnterStart: "ExitEnterStart",
	EvExitEnterCancel: "ExitEnterCancel",
	EvExitEnterFinished: "ExitEnterFinished",
	EvNewQuestGiverObject: "NewQuestGiverObject",
	EvFullQuestInfo: "FullQuestInfo",
	EvQuestProgressInfo: "QuestProgressInfo",
	EvQuestGiverInfoForPlayer: "QuestGiverInfoForPlayer",
	EvFullExpeditionInfo: "FullExpeditionInfo",
	EvExpeditionQuestProgressInfo: "ExpeditionQuestProgressInfo",
	EvInvitedToExpedition: "InvitedToExpedition",
	EvExpeditionRegistrationInfo: "ExpeditionRegistrationInfo",
	EvEnteringExpeditionStart: "EnteringExpeditionStart",
	EvEnteringExpeditionCancel: "EnteringExpeditionCancel",
	EvRewardGranted: "RewardGranted",
	EvArenaRegistrationInfo: "ArenaRegistrationInfo",
	EvEnteringArenaStart: "EnteringArenaStart",
	EvEnteringArenaCancel: "EnteringArenaCancel",
	EvEnteringArenaLockStart: "EnteringArenaLockStart",
	EvEnteringArenaLockCancel: "EnteringArenaLockCancel",
	EvInvitedToArenaMatch: "InvitedToArenaMatch",
	EvUsingHellgateShrine: "UsingHellgateShrine",
	EvEnteringHellgateLockStart: "EnteringHellgateLockStart",
	EvEnteringHellgateLockCancel: "EnteringHellgateLockCancel",
	EvPlayerCounts: "PlayerCounts",
	EvInCombatStateUpdate: "InCombatStateUpdate",
	EvOtherGrabbedLoot: "OtherGrabbedLoot",
	EvTreasureChestUsingStart: "TreasureChestUsingStart",
	EvTreasureChestUsingFinished: "TreasureChestUsingFinished",
	EvTreasureChestUsingCancel: "TreasureChestUsingCancel",
	EvTreasureChestUsingOpeningComplete: "TreasureChestUsingOpeningComplete",
	EvTreasureChestForceCloseInventory: "TreasureChestForceCloseInventory",
	EvLocalTreasuresUpdate: "LocalTreasuresUpdate",
	EvLootChestSpawnpointsUpdate: "LootChestSpawnpointsUpdate",
	EvPremiumChanged: "PremiumChanged",
	EvPremiumExtended: "PremiumExtended",
	EvPremiumLifeTimeRewardGained: "PremiumLifeTimeRewardGained",
	EvGoldPurchased: "GoldPurchased",
	EvLaborerGotUpgraded: "LaborerGotUpgraded",
	EvJournalGotFull: "JournalGotFull",
	EvJournalFillError: "JournalFillError",
	EvFriendRequest: "FriendRequest",
	EvFriendRequestInfos: "FriendRequestInfos",
	EvFriendInfos: "FriendInfos",
	EvFriendRequestAnswered: "FriendRequestAnswered",
	EvFriendOnlineStatus: "FriendOnlineStatus",
	EvFriendRequestCanceled: "FriendRequestCanceled",
	EvFriendRemoved: "FriendRemoved",
	EvFriendUpdated: "FriendUpdated",
	EvPartyLootItems: "PartyLootItems",
	EvPartyLootItemsRemoved: "PartyLootItemsRemoved",
	EvPartyLootItemTypesRemoved: "PartyLootItemTypesRemoved",
	EvReputationUpdate: "ReputationUpdate",
	EvDefenseUnitAttackBegin: "DefenseUnitAttackBegin",
	EvDefenseUnitAttackEnd: "DefenseUnitAttackEnd",
	EvDefenseUnitAttackDamage: "DefenseUnitAttackDamage",
	EvUnrestrictedPvpZoneUpdate: "UnrestrictedPvpZoneUpdate",
	EvUnrestrictedPvpZoneStatus: "UnrestrictedPvpZoneStatus",
	EvReputationImplicationUpdate: "ReputationImplicationUpdate",
	EvNewMountObject: "NewMountObject",
	EvMountHealthUpdate: "MountHealthUpdate",
	EvMountCooldownUpdate: "MountCooldownUpdate",
	EvNewExpeditionAgent: "NewExpeditionAgent",
	EvNewExpeditionCheckPoint: "NewExpeditionCheckPoint",
	EvExpeditionStartEvent: "ExpeditionStartEvent",
	EvVoteEvent: "VoteEvent",
	EvRatingEvent: "RatingEvent",
	EvNewArenaAgent: "NewArenaAgent",
	EvBoostFarmable: "BoostFarmable",
	EvUseFunction: "UseFunction",
	EvNewPortalEntrance: "NewPortalEntrance",
	EvNewPortalExit: "NewPortalExit",
	EvNewRandomDungeonExit: "NewRandomDungeonExit",
	EvWaitingQueueUpdate: "WaitingQueueUpdate",
	EvPlayerMovementRateUpdate: "PlayerMovementRateUpdate",
	EvObserveStart: "ObserveStart",
	EvMinimapZergs: "MinimapZergs",
	EvMinimapSmartClusterZergs: "MinimapSmartClusterZergs",
	EvPaymentTransactions: "PaymentTransactions",
	EvPerformanceStatsUpdate: "PerformanceStatsUpdate",
	EvOverloadModeUpdate: "OverloadModeUpdate",
	EvDebugDrawEvent: "DebugDrawEvent",
	EvRecordCameraMove: "RecordCameraMove",
	EvRecordStart: "RecordStart",
	EvDeliverCarriableObjectStart: "DeliverCarriableObjectStart",
	EvDeliverCarriableObjectCancel: "DeliverCarriableObjectCancel",
	EvDeliverCarriableObjectReset: "DeliverCarriableObjectReset",
	EvDeliverCarriableObjectFinished: "DeliverCarriableObjectFinished",
	EvTerritoryClaimStart: "TerritoryClaimStart",
	EvTerritoryClaimCancel: "TerritoryClaimCancel",
	EvTerritoryClaimFinished: "TerritoryClaimFinished",
	EvTerritoryScheduleResult: "TerritoryScheduleResult",
	EvTerritoryUpgradeWithPowerCrystalResult: "TerritoryUpgradeWithPowerCrystalResult",
	EvReceiveCarriableObjectStart: "ReceiveCarriableObjectStart",
	EvReceiveCarriableObjectFinished: "ReceiveCarriableObjectFinished",
	EvUpdateAccountState: "UpdateAccountState",
	EvStartDeterministicRoam: "StartDeterministicRoam",
	EvGuildFullAccessTagsUpdated: "GuildFullAccessTagsUpdated",
	EvGuildAccessTagUpdated: "GuildAccessTagUpdated",
	EvGvgSeasonUpdate: "GvgSeasonUpdate",
	EvGvgSeasonCheatCommand: "GvgSeasonCheatCommand",
	EvSeasonPointsByKillingBooster: "SeasonPointsByKillingBooster",
	EvFishingStart: "FishingStart",
	EvFishingCast: "FishingCast",
	EvFishingCatch: "FishingCatch",
	EvFishingFinished: "FishingFinished",
	EvFishingCancel: "FishingCancel",
	EvNewFloatObject: "NewFloatObject",
	EvNewFishingZoneObject: "NewFishingZoneObject",
	EvFishingMiniGame: "FishingMiniGame",
	EvAlbionJournalAchievementCompleted: "AlbionJournalAchievementCompleted",
	EvUpdatePuppet: "UpdatePuppet",
	EvChangeFlaggingFinished: "ChangeFlaggingFinished",
	EvNewOutpostObject: "NewOutpostObject",
	EvOutpostUpdate: "OutpostUpdate",
	EvOutpostClaimed: "OutpostClaimed",
	EvOverChargeEnd: "OverChargeEnd",
	EvOverChargeStatus: "OverChargeStatus",
	EvPartyFinderFullUpdate: "PartyFinderFullUpdate",
	EvPartyFinderUpdate: "PartyFinderUpdate",
	EvPartyFinderApplicantsUpdate: "PartyFinderApplicantsUpdate",
	EvPartyFinderEquipmentSnapshot: "PartyFinderEquipmentSnapshot",
	EvPartyFinderJoinRequestDeclined: "PartyFinderJoinRequestDeclined",
	EvNewUnlockedPersonalSeasonRewards: "NewUnlockedPersonalSeasonRewards",
	EvPersonalSeasonPointsGained: "PersonalSeasonPointsGained",
	EvPersonalSeasonPastSeasonDataEvent: "PersonalSeasonPastSeasonDataEvent",
	EvMatchLootChestOpeningStart: "MatchLootChestOpeningStart",
	EvMatchLootChestOpeningFinished: "MatchLootChestOpeningFinished",
	EvMatchLootChestOpeningCancel: "MatchLootChestOpeningCancel",
	EvNotifyCrystalMatchReward: "NotifyCrystalMatchReward",
	EvCrystalRealmFeedback: "CrystalRealmFeedback",
	EvNewLocationMarker: "NewLocationMarker",
	EvNewTutorialBlocker: "NewTutorialBlocker",
	EvNewTileSwitch: "NewTileSwitch",
	EvNewInformationProvider: "NewInformationProvider",
	EvNewDynamicGuildLogo: "NewDynamicGuildLogo",
	EvNewDecoration: "NewDecoration",
	EvTutorialUpdate: "TutorialUpdate",
	EvTriggerHintBox: "TriggerHintBox",
	EvRandomDungeonPositionInfo: "RandomDungeonPositionInfo",
	EvNewLootChest: "NewLootChest",
	EvUpdateLootChest: "UpdateLootChest",
	EvLootChestOpened: "LootChestOpened",
	EvUpdateLootProtectedByMobsWithMinimapDisplay: "UpdateLootProtectedByMobsWithMinimapDisplay",
	EvNewShrine: "NewShrine",
	EvUpdateShrine: "UpdateShrine",
	EvUpdateRoom: "UpdateRoom",
	EvNewMobSoul: "NewMobSoul",
	EvNewHellgateShrine: "NewHellgateShrine",
	EvUpdateHellgateShrine: "UpdateHellgateShrine",
	EvActivateHellgateExit: "ActivateHellgateExit",
	EvMutePlayerUpdate: "MutePlayerUpdate",
	EvShopTileUpdate: "ShopTileUpdate",
	EvShopUpdate: "ShopUpdate",
	EvAntiCheatKick: "AntiCheatKick",
	EvBattlEyeServerMessage: "BattlEyeServerMessage",
	EvUnlockVanityUnlock: "UnlockVanityUnlock",
	EvAvatarUnlocked: "AvatarUnlocked",
	EvCustomizationChanged: "CustomizationChanged",
	EvBaseVaultInfo: "BaseVaultInfo",
	EvGuildVaultInfo: "GuildVaultInfo",
	EvBankVaultInfo: "BankVaultInfo",
	EvRecoveryVaultPlayerInfo: "RecoveryVaultPlayerInfo",
	EvRecoveryVaultGuildInfo: "RecoveryVaultGuildInfo",
	EvUpdateWardrobe: "UpdateWardrobe",
	EvCastlePhaseChanged: "CastlePhaseChanged",
	EvGuildAccountLogEvent: "GuildAccountLogEvent",
	EvNewHideoutObject: "NewHideoutObject",
	EvNewHideoutManagement: "NewHideoutManagement",
	EvNewHideoutExit: "NewHideoutExit",
	EvInitHideoutAttackStart: "InitHideoutAttackStart",
	EvInitHideoutAttackCancel: "InitHideoutAttackCancel",
	EvInitHideoutAttackFinished: "InitHideoutAttackFinished",
	EvHideoutManagementUpdate: "HideoutManagementUpdate",
	EvHideoutUpgradeWithPowerCrystalResult: "HideoutUpgradeWithPowerCrystalResult",
	EvIpChanged: "IpChanged",
	EvSmartClusterQueueUpdateInfo: "SmartClusterQueueUpdateInfo",
	EvSmartClusterQueueActiveInfo: "SmartClusterQueueActiveInfo",
	EvSmartClusterQueueKickWarning: "SmartClusterQueueKickWarning",
	EvSmartClusterQueueInvite: "SmartClusterQueueInvite",
	EvReceivedGvgSeasonPoints: "ReceivedGvgSeasonPoints",
	EvTowerPowerPointUpdate: "TowerPowerPointUpdate",
	EvOpenWorldAttackScheduleStart: "OpenWorldAttackScheduleStart",
	EvOpenWorldAttackScheduleFinished: "OpenWorldAttackScheduleFinished",
	EvOpenWorldAttackScheduleCancel: "OpenWorldAttackScheduleCancel",
	EvOpenWorldAttackConquerStart: "OpenWorldAttackConquerStart",
	EvOpenWorldAttackConquerFinished: "OpenWorldAttackConquerFinished",
	EvOpenWorldAttackConquerCancel: "OpenWorldAttackConquerCancel",
	EvOpenWorldAttackConquerStatus: "OpenWorldAttackConquerStatus",
	EvOpenWorldAttackStart: "OpenWorldAttackStart",
	EvOpenWorldAttackEnd: "OpenWorldAttackEnd",
	EvNewRandomResourceBlocker: "NewRandomResourceBlocker",
	EvNewHomeObject: "NewHomeObject",
	EvHideoutObjectUpdate: "HideoutObjectUpdate",
	EvUpdateInfamy: "UpdateInfamy",
	EvMinimapPositionMarkers: "MinimapPositionMarkers",
	EvNewTunnelExit: "NewTunnelExit",
	EvCorruptedDungeonUpdate: "CorruptedDungeonUpdate",
	EvCorruptedDungeonStatus: "CorruptedDungeonStatus",
	EvCorruptedDungeonInfamy: "CorruptedDungeonInfamy",
	EvHellgateRestrictedAreaUpdate: "HellgateRestrictedAreaUpdate",
	EvHellgateInfamy: "HellgateInfamy",
	EvHellgateStatus: "HellgateStatus",
	EvHellgateStatusUpdate: "HellgateStatusUpdate",
	EvHellgateSuspense: "HellgateSuspense",
	EvReplaceSpellSlotWithMultiSpell: "ReplaceSpellSlotWithMultiSpell",
	EvNewCorruptedShrine: "NewCorruptedShrine",
	EvUpdateCorruptedShrine: "UpdateCorruptedShrine",
	EvCorruptedShrineUsageStart: "CorruptedShrineUsageStart",
	EvCorruptedShrineUsageCancel: "CorruptedShrineUsageCancel",
	EvExitUsed: "ExitUsed",
	EvLinkedToObject: "LinkedToObject",
	EvLinkToObjectBroken: "LinkToObjectBroken",
	EvEstimatedMarketValueUpdate: "EstimatedMarketValueUpdate",
	EvStuckCancel: "StuckCancel",
	EvDungonEscapeReady: "DungonEscapeReady",
	EvFactionWarfareClusterState: "FactionWarfareClusterState",
	EvFactionWarfareHasUnclaimedWeeklyReportsEvent: "FactionWarfareHasUnclaimedWeeklyReportsEvent",
	EvSimpleFeedback: "SimpleFeedback",
	EvSmartClusterQueueSkipClusterError: "SmartClusterQueueSkipClusterError",
	EvXignCodeEvent: "XignCodeEvent",
	EvBatchUseItemStart: "BatchUseItemStart",
	EvBatchUseItemEnd: "BatchUseItemEnd",
	EvRedZonePlayerNotification: "RedZonePlayerNotification",
	EvRedZoneEventCheatCleanup: "RedZoneEventCheatCleanup",
	EvRedZoneFortressEventChestOpened: "RedZoneFortressEventChestOpened",
	EvRedZoneWorldMapEvent: "RedZoneWorldMapEvent",
	EvFactionWarfareStats: "FactionWarfareStats",
	EvUpdateFactionBalanceFactors: "UpdateFactionBalanceFactors",
	EvFactionEnlistmentChanged: "FactionEnlistmentChanged",
	EvUpdateFactionRank: "UpdateFactionRank",
	EvFactionWarfareCampaignRewardsUnlocked: "FactionWarfareCampaignRewardsUnlocked",
	EvFeaturedFeatureUpdate: "FeaturedFeatureUpdate",
	EvNewCarriableObject: "NewCarriableObject",
	EvMinimapCrystalPositionMarker: "MinimapCrystalPositionMarker",
	EvCarriedObjectUpdate: "CarriedObjectUpdate",
	EvPickupCarriableObjectStart: "PickupCarriableObjectStart",
	EvPickupCarriableObjectCancel: "PickupCarriableObjectCancel",
	EvPickupCarriableObjectFinished: "PickupCarriableObjectFinished",
	EvDoSimpleActionStart: "DoSimpleActionStart",
	EvDoSimpleActionCancel: "DoSimpleActionCancel",
	EvDoSimpleActionFinished: "DoSimpleActionFinished",
	EvNotifyGuestAccountVerified: "NotifyGuestAccountVerified",
	EvMightAndFavorReceivedEvent: "MightAndFavorReceivedEvent",
	EvWeeklyPvpChallengeRewardStateUpdate: "WeeklyPvpChallengeRewardStateUpdate",
	EvNewUnlockedPvpSeasonChallengeRewards: "NewUnlockedPvpSeasonChallengeRewards",
	EvStaticDungeonEntrancesDungeonEventStatusUpdates: "StaticDungeonEntrancesDungeonEventStatusUpdates",
	EvStaticDungeonDungeonValueUpdate: "StaticDungeonDungeonValueUpdate",
	EvStaticDungeonEntranceDungeonEventsAborted: "StaticDungeonEntranceDungeonEventsAborted",
	EvInAppPurchaseConfirmedGooglePlay: "InAppPurchaseConfirmedGooglePlay",
	EvFeatureSwitchInfo: "FeatureSwitchInfo",
	EvPartyJoinRequestAborted: "PartyJoinRequestAborted",
	EvPartyInviteAborted: "PartyInviteAborted",
	EvPartyStartHuntRequest: "PartyStartHuntRequest",
	EvPartyStartHuntRequested: "PartyStartHuntRequested",
	EvPartyStartHuntRequestAnswer: "PartyStartHuntRequestAnswer",
	EvPartyPlayerLeaveScheduled: "PartyPlayerLeaveScheduled",
	EvGuildInviteDeclined: "GuildInviteDeclined",
	EvCancelMultiSpellSlots: "CancelMultiSpellSlots",
	EvNewVisualEventObject: "NewVisualEventObject",
	EvCastleClaimProgress: "CastleClaimProgress",
	EvCastleClaimProgressLogo: "CastleClaimProgressLogo",
	EvTownPortalUpdateState: "TownPortalUpdateState",
	EvTownPortalFailed: "TownPortalFailed",
	EvConsumableVanityChargesAdded: "ConsumableVanityChargesAdded",
	EvFestivitiesUpdate: "FestivitiesUpdate",
	EvNewBannerObject: "NewBannerObject",
	EvNewMistsImmediateReturnExit: "NewMistsImmediateReturnExit",
	EvMistsPlayerJoinedInfo: "MistsPlayerJoinedInfo",
	EvNewMistsStaticEntrance: "NewMistsStaticEntrance",
	EvNewMistsOpenWorldExit: "NewMistsOpenWorldExit",
	EvNewTunnelExitTemp: "NewTunnelExitTemp",
	EvNewMistsWispSpawn: "NewMistsWispSpawn",
	EvMistsWispSpawnStateChange: "MistsWispSpawnStateChange",
	EvNewMistsCityEntrance: "NewMistsCityEntrance",
	EvNewMistsCityRoadsEntrance: "NewMistsCityRoadsEntrance",
	EvMistsCityRoadsEntrancePartyStateUpdate: "MistsCityRoadsEntrancePartyStateUpdate",
	EvMistsCityRoadsEntranceClearStateForParty: "MistsCityRoadsEntranceClearStateForParty",
	EvMistsEntranceDataChanged: "MistsEntranceDataChanged",
	EvNewCagedObject: "NewCagedObject",
	EvCagedObjectStateUpdated: "CagedObjectStateUpdated",
	EvEntrancePartyBindingCreated: "EntrancePartyBindingCreated",
	EvEntrancePartyBindingCleared: "EntrancePartyBindingCleared",
	EvEntrancePartyBindingInfos: "EntrancePartyBindingInfos",
	EvNewMistsBorderExit: "NewMistsBorderExit",
	EvNewMistsDungeonExit: "NewMistsDungeonExit",
	EvLocalQuestInfos: "LocalQuestInfos",
	EvLocalQuestStarted: "LocalQuestStarted",
	EvLocalQuestActive: "LocalQuestActive",
	EvLocalQuestInactive: "LocalQuestInactive",
	EvLocalQuestProgressUpdate: "LocalQuestProgressUpdate",
	EvNewUnrestrictedPvpZone: "NewUnrestrictedPvpZone",
	EvTemporaryFlaggingStatusUpdate: "TemporaryFlaggingStatusUpdate",
	EvSpellTestPerformanceUpdate: "SpellTestPerformanceUpdate",
	EvTransformation: "Transformation",
	EvTransformationEnd: "TransformationEnd",
	EvUpdateTrustlevel: "UpdateTrustlevel",
	EvRevealHiddenTimeStamps: "RevealHiddenTimeStamps",
	EvModifyItemTraitFinished: "ModifyItemTraitFinished",
	EvRerollItemTraitValueFinished: "RerollItemTraitValueFinished",
	EvHuntQuestProgressInfo: "HuntQuestProgressInfo",
	EvHuntStarted: "HuntStarted",
	EvHuntFinished: "HuntFinished",
	EvHuntAborted: "HuntAborted",
	EvHuntMissionStepStateUpdate: "HuntMissionStepStateUpdate",
	EvNewHuntTrack: "NewHuntTrack",
	EvHuntMissionUpdate: "HuntMissionUpdate",
	EvHuntQuestMissionProgressUpdate: "HuntQuestMissionProgressUpdate",
	EvHuntTrackUsed: "HuntTrackUsed",
	EvHuntTrackUseableAgain: "HuntTrackUseableAgain",
	EvMinimapHuntTrackMarkers: "MinimapHuntTrackMarkers",
	EvNoTracksFound: "NoTracksFound",
	EvHuntQuestAborted: "HuntQuestAborted",
	EvInteractWithTrackStart: "InteractWithTrackStart",
	EvInteractWithTrackCancel: "InteractWithTrackCancel",
	EvInteractWithTrackFinished: "InteractWithTrackFinished",
	EvNewDynamicCompound: "NewDynamicCompound",
	EvLegendaryItemDestroyed: "LegendaryItemDestroyed",
	EvAttunementInfo: "AttunementInfo",
	EvTerritoryClaimRaidedRawEnergyCrystalResult: "TerritoryClaimRaidedRawEnergyCrystalResult",
	EvCarriedObjectExpiryWarning: "CarriedObjectExpiryWarning",
	EvCarriedObjectExpired: "CarriedObjectExpired",
	EvTerritoryRaidStart: "TerritoryRaidStart",
	EvTerritoryRaidCancel: "TerritoryRaidCancel",
	EvTerritoryRaidFinished: "TerritoryRaidFinished",
	EvTerritoryRaidResult: "TerritoryRaidResult",
	EvTerritoryMonolithActiveRaidStatus: "TerritoryMonolithActiveRaidStatus",
	EvTerritoryMonolithActiveRaidCancelled: "TerritoryMonolithActiveRaidCancelled",
	EvMonolithEnergyStorageUpdate: "MonolithEnergyStorageUpdate",
	EvMonolithNextScheduledOpenWorldAttackUpdate: "MonolithNextScheduledOpenWorldAttackUpdate",
	EvMonolithProtectedBuildingsDamageReductionUpdate: "MonolithProtectedBuildingsDamageReductionUpdate",
	EvNewBuildingBaseEvent: "NewBuildingBaseEvent",
	EvNewFortificationBuilding: "NewFortificationBuilding",
	EvNewCastleGateBuilding: "NewCastleGateBuilding",
	EvBuildingDurabilityUpdate: "BuildingDurabilityUpdate",
	EvMonolithFortificationPointsUpdate: "MonolithFortificationPointsUpdate",
	EvFortificationBuildingUpgradeInfo: "FortificationBuildingUpgradeInfo",
	EvFortificationBuildingsDamageStateUpdate: "FortificationBuildingsDamageStateUpdate",
	EvSiegeNotificationEvent: "SiegeNotificationEvent",
	EvUpdateEnemyWarBannerActive: "UpdateEnemyWarBannerActive",
	EvTerritoryAnnouncePlayerEjection: "TerritoryAnnouncePlayerEjection",
	EvCastleGateSwitchUseStarted: "CastleGateSwitchUseStarted",
	EvCastleGateSwitchUseFinished: "CastleGateSwitchUseFinished",
	EvFortificationBuildingWillDowngrade: "FortificationBuildingWillDowngrade",
	EvBotCommand: "BotCommand",
	EvJournalAchievementProgressUpdate: "JournalAchievementProgressUpdate",
	EvJournalClaimableRewardUpdate: "JournalClaimableRewardUpdate",
	EvKeySync: "KeySync",
	EvLocalQuestAreaGone: "LocalQuestAreaGone",
	EvDynamicTemplate: "DynamicTemplate",
	EvDynamicTemplateForcedStateChange: "DynamicTemplateForcedStateChange",
	EvNewOutlandsTeleportationPortal: "NewOutlandsTeleportationPortal",
	EvNewOutlandsTeleportationReturnPortal: "NewOutlandsTeleportationReturnPortal",
	EvOutlandsTeleportationBindingCleared: "OutlandsTeleportationBindingCleared",
	EvOutlandsTeleportationReturnPortalUpdateEvent: "OutlandsTeleportationReturnPortalUpdateEvent",
	EvPlayerUsedOutlandsTeleportationPortal: "PlayerUsedOutlandsTeleportationPortal",
	EvEncumberedRestricted: "EncumberedRestricted",
	EvNewPiledObject: "NewPiledObject",
	EvPiledObjectStateChanged: "PiledObjectStateChanged",
	EvNewSmugglerCrateDeliveryStation: "NewSmugglerCrateDeliveryStation",
	EvKillRewardedNoFame: "KillRewardedNoFame",
	EvPickupFromPiledObjectStart: "PickupFromPiledObjectStart",
	EvPickupFromPiledObjectCancel: "PickupFromPiledObjectCancel",
	EvPickupFromPiledObjectReset: "PickupFromPiledObjectReset",
	EvPickupFromPiledObjectFinished: "PickupFromPiledObjectFinished",
	EvArmoryActivityChange: "ArmoryActivityChange",
	EvNewKillTrophyFurnitureBuilding: "NewKillTrophyFurnitureBuilding",
	EvHellDungeonsPlayerJoinedInfo: "HellDungeonsPlayerJoinedInfo",
	EvNewTileSwitchTrigger: "NewTileSwitchTrigger",
	EvNewMultiRewardObject: "NewMultiRewardObject",
	EvNewHellDungeonSoulShrineObject: "NewHellDungeonSoulShrineObject",
	EvHellDungeonSoulShrineStateUpdate: "HellDungeonSoulShrineStateUpdate",
	EvNewResurrectionShrine: "NewResurrectionShrine",
	EvUpdateResurrectionShrine: "UpdateResurrectionShrine",
	EvStandTimeFinished: "StandTimeFinished",
	EvEpicAchievementAndStatsUpdate: "EpicAchievementAndStatsUpdate",
	EvSpectateTargetAfterDeathUpdate: "SpectateTargetAfterDeathUpdate",
	EvSpectateTargetAfterDeathEnded: "SpectateTargetAfterDeathEnded",
	EvNewHellDungeonUpwardExit: "NewHellDungeonUpwardExit",
	EvNewHellDungeonSoulExit: "NewHellDungeonSoulExit",
	EvNewHellDungeonDownwardExit: "NewHellDungeonDownwardExit",
	EvNewHellDungeonChestExit: "NewHellDungeonChestExit",
	EvNewCorruptedStaticEntrance: "NewCorruptedStaticEntrance",
	EvNewHellDungeonStaticEntrance: "NewHellDungeonStaticEntrance",
	EvUpdateHellDungeonStaticEntranceState: "UpdateHellDungeonStaticEntranceState",
	EvDebugTriggerHellDungeonShutdownStart: "DebugTriggerHellDungeonShutdownStart",
	EvFullJournalQuestInfo: "FullJournalQuestInfo",
	EvJournalQuestProgressInfo: "JournalQuestProgressInfo",
	EvNewHellDungeonRoomShrineObject: "NewHellDungeonRoomShrineObject",
	EvHellDungeonRoomShrineStateUpdate: "HellDungeonRoomShrineStateUpdate",
	EvSimpleBehaviourBuildingStateUpdate: "SimpleBehaviourBuildingStateUpdate",
	EvSetTimeScaling: "SetTimeScaling",
	EvStopTimeScaling: "StopTimeScaling",
	EvKeyValidation: "KeyValidation",
	EvPlayerJoinMapMarkerTimerStates: "PlayerJoinMapMarkerTimerStates",
	EvNewMapMarkerTimer: "NewMapMarkerTimer",
	EvRemoveMapMarkerTimer: "RemoveMapMarkerTimer",
	EvNewFactionFortressObject: "NewFactionFortressObject",
	EvFactionFortressAnnouncePlayerEjection: "FactionFortressAnnouncePlayerEjection",
	EvRewardFactionWarfareSupply: "RewardFactionWarfareSupply",
	EvFactionCaptureAreaProgressUpdate: "FactionCaptureAreaProgressUpdate",
	EvFactionFortressClaimed: "FactionFortressClaimed",
	EvFactionFortressWeaponCachesSpawned: "FactionFortressWeaponCachesSpawned",
	EvFactionFortressWeaponCacheClaimed: "FactionFortressWeaponCacheClaimed",
	EvFactionFortressFightStateUpdate: "FactionFortressFightStateUpdate",
	EvFactionFortressCutoffFightStateUpdate: "FactionFortressCutoffFightStateUpdate",
	EvFactionFortressFightEnded: "FactionFortressFightEnded",
	EvNewFactionWarfarePortal: "NewFactionWarfarePortal",
	EvFactionPortalTargetUpdate: "FactionPortalTargetUpdate",
	EvFactionFortressFightStartedInRemoteClusterEvent: "FactionFortressFightStartedInRemoteClusterEvent",
	EvFactionFortressFightFinishedInRemoteClusterEvent: "FactionFortressFightFinishedInRemoteClusterEvent",
	EvFactionDuchySupplyWarDefensiveVictoryEvent: "FactionDuchySupplyWarDefensiveVictoryEvent",
	EvFactionDuchyReconnectedFromCutoffEvent: "FactionDuchyReconnectedFromCutoffEvent",
	EvFactionFortressCutoffFightCancelledByClusterOwnerChangeEvent: "FactionFortressCutoffFightCancelledByClusterOwnerChangeEvent",
	EvFactionDuchyEnteredCutoffStateEvent: "FactionDuchyEnteredCutoffStateEvent",
	EvLeaveProtectionStateUpdate: "LeaveProtectionStateUpdate",
	EvRedZoneEventStandings: "RedZoneEventStandings",
	EvNewFactionBattleStandardDeliveryStation: "NewFactionBattleStandardDeliveryStation",
	EvNewLoreSnippetObject: "NewLoreSnippetObject",
	EvLoreSnippetObjectStateUpdate: "LoreSnippetObjectStateUpdate",
	EvLoreSnippedClaimed: "LoreSnippedClaimed",
	EvLoreSnippetStatesChangedByCheat: "LoreSnippetStatesChangedByCheat",
	EvNewTeleporterNode: "NewTeleporterNode",
	EvTeleporterNodeStateChanged: "TeleporterNodeStateChanged",
	EvTeleporterConnectionsFullStateUpdate: "TeleporterConnectionsFullStateUpdate",
	EvTeleporterConnectionStateChanged: "TeleporterConnectionStateChanged",
	EvRetrieveCarriableObjectStart: "RetrieveCarriableObjectStart",
	EvRetrieveCarriableObjectCancel: "RetrieveCarriableObjectCancel",
	EvRetrieveCarriableObjectReset: "RetrieveCarriableObjectReset",
	EvRetrieveCarriableObjectFinished: "RetrieveCarriableObjectFinished",
	EvLosingCarriableObjectStart: "LosingCarriableObjectStart",
	EvLosingCarriableObjectFinished: "LosingCarriableObjectFinished",
}

func (c EventCode) String() string { if n,ok:=eventCodeNames[c]; ok { return n }; return "Event("+itoa(int(c))+")" }
func IsKnownEventCode(c EventCode) bool { _,ok:=eventCodeNames[c]; return ok }

