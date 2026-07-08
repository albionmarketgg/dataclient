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
	EvSpellCooldownUpdate EventCode = 250
	EvNewHellgateExitPortal EventCode = 251
	EvNewExpeditionExit EventCode = 252
	EvNewExpeditionNarrator EventCode = 253
	EvExitEnterStart EventCode = 254
	EvExitEnterCancel EventCode = 255
	EvExitEnterFinished EventCode = 256
	EvNewQuestGiverObject EventCode = 257
	EvFullQuestInfo EventCode = 258
	EvQuestProgressInfo EventCode = 259
	EvQuestGiverInfoForPlayer EventCode = 260
	EvFullExpeditionInfo EventCode = 261
	EvExpeditionQuestProgressInfo EventCode = 262
	EvInvitedToExpedition EventCode = 263
	EvExpeditionRegistrationInfo EventCode = 264
	EvEnteringExpeditionStart EventCode = 265
	EvEnteringExpeditionCancel EventCode = 266
	EvRewardGranted EventCode = 267
	EvArenaRegistrationInfo EventCode = 268
	EvEnteringArenaStart EventCode = 269
	EvEnteringArenaCancel EventCode = 270
	EvEnteringArenaLockStart EventCode = 271
	EvEnteringArenaLockCancel EventCode = 272
	EvInvitedToArenaMatch EventCode = 273
	EvUsingHellgateShrine EventCode = 274
	EvEnteringHellgateLockStart EventCode = 275
	EvEnteringHellgateLockCancel EventCode = 276
	EvPlayerCounts EventCode = 277
	EvInCombatStateUpdate EventCode = 278
	EvOtherGrabbedLoot EventCode = 279
	EvTreasureChestUsingStart EventCode = 280
	EvTreasureChestUsingFinished EventCode = 281
	EvTreasureChestUsingCancel EventCode = 282
	EvTreasureChestUsingOpeningComplete EventCode = 283
	EvTreasureChestForceCloseInventory EventCode = 284
	EvLocalTreasuresUpdate EventCode = 285
	EvLootChestSpawnpointsUpdate EventCode = 286
	EvPremiumChanged EventCode = 287
	EvPremiumExtended EventCode = 288
	EvPremiumLifeTimeRewardGained EventCode = 289
	EvGoldPurchased EventCode = 290
	EvLaborerGotUpgraded EventCode = 291
	EvJournalGotFull EventCode = 292
	EvJournalFillError EventCode = 293
	EvFriendRequest EventCode = 294
	EvFriendRequestInfos EventCode = 295
	EvFriendInfos EventCode = 296
	EvFriendRequestAnswered EventCode = 297
	EvFriendOnlineStatus EventCode = 298
	EvFriendRequestCanceled EventCode = 299
	EvFriendRemoved EventCode = 300
	EvFriendUpdated EventCode = 301
	EvPartyLootItems EventCode = 302
	EvPartyLootItemsRemoved EventCode = 303
	EvPartyLootItemTypesRemoved EventCode = 304
	EvReputationUpdate EventCode = 305
	EvDefenseUnitAttackBegin EventCode = 306
	EvDefenseUnitAttackEnd EventCode = 307
	EvDefenseUnitAttackDamage EventCode = 308
	EvUnrestrictedPvpZoneUpdate EventCode = 309
	EvUnrestrictedPvpZoneStatus EventCode = 310
	EvReputationImplicationUpdate EventCode = 311
	EvNewMountObject EventCode = 312
	EvMountHealthUpdate EventCode = 313
	EvMountCooldownUpdate EventCode = 314
	EvNewExpeditionAgent EventCode = 315
	EvNewExpeditionCheckPoint EventCode = 316
	EvExpeditionStartEvent EventCode = 317
	EvVoteEvent EventCode = 318
	EvRatingEvent EventCode = 319
	EvNewArenaAgent EventCode = 320
	EvBoostFarmable EventCode = 321
	EvUseFunction EventCode = 322
	EvNewPortalEntrance EventCode = 323
	EvNewPortalExit EventCode = 324
	EvNewRandomDungeonExit EventCode = 325
	EvWaitingQueueUpdate EventCode = 326
	EvPlayerMovementRateUpdate EventCode = 327
	EvObserveStart EventCode = 328
	EvMinimapZergs EventCode = 329
	EvMinimapSmartClusterZergs EventCode = 330
	EvPaymentTransactions EventCode = 331
	EvPerformanceStatsUpdate EventCode = 332
	EvOverloadModeUpdate EventCode = 333
	EvDebugDrawEvent EventCode = 334
	EvRecordCameraMove EventCode = 335
	EvRecordStart EventCode = 336
	EvDeliverCarriableObjectStart EventCode = 337
	EvDeliverCarriableObjectCancel EventCode = 338
	EvDeliverCarriableObjectReset EventCode = 339
	EvDeliverCarriableObjectFinished EventCode = 340
	EvTerritoryClaimStart EventCode = 341
	EvTerritoryClaimCancel EventCode = 342
	EvTerritoryClaimFinished EventCode = 343
	EvTerritoryScheduleResult EventCode = 344
	EvTerritoryUpgradeWithPowerCrystalResult EventCode = 345
	EvReceiveCarriableObjectStart EventCode = 346
	EvReceiveCarriableObjectFinished EventCode = 347
	EvUpdateAccountState EventCode = 348
	EvStartDeterministicRoam EventCode = 349
	EvGuildFullAccessTagsUpdated EventCode = 350
	EvGuildAccessTagUpdated EventCode = 351
	EvGvgSeasonUpdate EventCode = 352
	EvGvgSeasonCheatCommand EventCode = 353
	EvSeasonPointsByKillingBooster EventCode = 354
	EvFishingStart EventCode = 355
	EvFishingCast EventCode = 356
	EvFishingCatch EventCode = 357
	EvFishingFinished EventCode = 358
	EvFishingCancel EventCode = 359
	EvNewFloatObject EventCode = 360
	EvNewFishingZoneObject EventCode = 361
	EvFishingMiniGame EventCode = 362
	EvAlbionJournalAchievementCompleted EventCode = 363
	EvUpdatePuppet EventCode = 364
	EvChangeFlaggingFinished EventCode = 365
	EvNewOutpostObject EventCode = 366
	EvOutpostUpdate EventCode = 367
	EvOutpostClaimed EventCode = 368
	EvOverChargeEnd EventCode = 369
	EvOverChargeStatus EventCode = 370
	EvPartyFinderFullUpdate EventCode = 371
	EvPartyFinderUpdate EventCode = 372
	EvPartyFinderApplicantsUpdate EventCode = 373
	EvPartyFinderEquipmentSnapshot EventCode = 374
	EvPartyFinderJoinRequestDeclined EventCode = 375
	EvNewUnlockedPersonalSeasonRewards EventCode = 376
	EvPersonalSeasonPointsGained EventCode = 377
	EvPersonalSeasonPastSeasonDataEvent EventCode = 378
	EvMatchLootChestOpeningStart EventCode = 379
	EvMatchLootChestOpeningFinished EventCode = 380
	EvMatchLootChestOpeningCancel EventCode = 381
	EvNotifyCrystalMatchReward EventCode = 382
	EvCrystalRealmFeedback EventCode = 383
	EvNewLocationMarker EventCode = 384
	EvNewTutorialBlocker EventCode = 385
	EvNewTileSwitch EventCode = 386
	EvNewInformationProvider EventCode = 387
	EvNewDynamicGuildLogo EventCode = 388
	EvNewDecoration EventCode = 389
	EvTutorialUpdate EventCode = 390
	EvTriggerHintBox EventCode = 391
	EvRandomDungeonPositionInfo EventCode = 392
	EvNewLootChest EventCode = 393
	EvUpdateLootChest EventCode = 394
	EvLootChestOpened EventCode = 395
	EvUpdateLootProtectedByMobsWithMinimapDisplay EventCode = 396
	EvNewShrine EventCode = 397
	EvUpdateShrine EventCode = 398
	EvUpdateRoom EventCode = 399
	EvNewMobSoul EventCode = 400
	EvNewHellgateShrine EventCode = 401
	EvUpdateHellgateShrine EventCode = 402
	EvActivateHellgateExit EventCode = 403
	EvMutePlayerUpdate EventCode = 404
	EvShopTileUpdate EventCode = 405
	EvShopUpdate EventCode = 406
	EvAntiCheatKick EventCode = 407
	EvBattlEyeServerMessage EventCode = 408
	EvUnlockVanityUnlock EventCode = 409
	EvAvatarUnlocked EventCode = 410
	EvCustomizationChanged EventCode = 411
	EvBaseVaultInfo EventCode = 412
	EvGuildVaultInfo EventCode = 413
	EvBankVaultInfo EventCode = 414
	EvRecoveryVaultPlayerInfo EventCode = 415
	EvRecoveryVaultGuildInfo EventCode = 416
	EvUpdateWardrobe EventCode = 417
	EvCastlePhaseChanged EventCode = 418
	EvGuildAccountLogEvent EventCode = 419
	EvNewHideoutObject EventCode = 420
	EvNewHideoutManagement EventCode = 421
	EvNewHideoutExit EventCode = 422
	EvInitHideoutAttackStart EventCode = 423
	EvInitHideoutAttackCancel EventCode = 424
	EvInitHideoutAttackFinished EventCode = 425
	EvHideoutManagementUpdate EventCode = 426
	EvHideoutUpgradeWithPowerCrystalResult EventCode = 427
	EvIpChanged EventCode = 428
	EvSmartClusterQueueUpdateInfo EventCode = 429
	EvSmartClusterQueueActiveInfo EventCode = 430
	EvSmartClusterQueueKickWarning EventCode = 431
	EvSmartClusterQueueInvite EventCode = 432
	EvReceivedGvgSeasonPoints EventCode = 433
	EvTowerPowerPointUpdate EventCode = 434
	EvOpenWorldAttackScheduleStart EventCode = 435
	EvOpenWorldAttackScheduleFinished EventCode = 436
	EvOpenWorldAttackScheduleCancel EventCode = 437
	EvOpenWorldAttackConquerStart EventCode = 438
	EvOpenWorldAttackConquerFinished EventCode = 439
	EvOpenWorldAttackConquerCancel EventCode = 440
	EvOpenWorldAttackConquerStatus EventCode = 441
	EvOpenWorldAttackStart EventCode = 442
	EvOpenWorldAttackEnd EventCode = 443
	EvNewRandomResourceBlocker EventCode = 444
	EvNewHomeObject EventCode = 445
	EvHideoutObjectUpdate EventCode = 446
	EvUpdateInfamy EventCode = 447
	EvMinimapPositionMarkers EventCode = 448
	EvNewTunnelExit EventCode = 449
	EvCorruptedDungeonUpdate EventCode = 450
	EvCorruptedDungeonStatus EventCode = 451
	EvCorruptedDungeonInfamy EventCode = 452
	EvHellgateRestrictedAreaUpdate EventCode = 453
	EvHellgateInfamy EventCode = 454
	EvHellgateStatus EventCode = 455
	EvHellgateStatusUpdate EventCode = 456
	EvHellgateSuspense EventCode = 457
	EvReplaceSpellSlotWithMultiSpell EventCode = 458
	EvNewCorruptedShrine EventCode = 459
	EvUpdateCorruptedShrine EventCode = 460
	EvCorruptedShrineUsageStart EventCode = 461
	EvCorruptedShrineUsageCancel EventCode = 462
	EvExitUsed EventCode = 463
	EvLinkedToObject EventCode = 464
	EvLinkToObjectBroken EventCode = 465
	EvEstimatedMarketValueUpdate EventCode = 466
	EvStuckCancel EventCode = 467
	EvDungonEscapeReady EventCode = 468
	EvFactionWarfareClusterState EventCode = 469
	EvFactionWarfareHasUnclaimedWeeklyReportsEvent EventCode = 470
	EvSimpleFeedback EventCode = 471
	EvSmartClusterQueueSkipClusterError EventCode = 472
	EvXignCodeEvent EventCode = 473
	EvBatchUseItemStart EventCode = 474
	EvBatchUseItemEnd EventCode = 475
	EvRedZonePlayerNotification EventCode = 476
	EvRedZoneEventCheatCleanup EventCode = 477
	EvRedZoneFortressEventChestOpened EventCode = 478
	EvRedZoneWorldMapEvent EventCode = 479
	EvFactionWarfareStats EventCode = 480
	EvUpdateFactionBalanceFactors EventCode = 481
	EvFactionEnlistmentChanged EventCode = 482
	EvUpdateFactionRank EventCode = 483
	EvFactionWarfareCampaignRewardsUnlocked EventCode = 484
	EvFeaturedFeatureUpdate EventCode = 485
	EvNewCarriableObject EventCode = 486
	EvMinimapCrystalPositionMarker EventCode = 487
	EvCarriedObjectUpdate EventCode = 488
	EvPickupCarriableObjectStart EventCode = 489
	EvPickupCarriableObjectCancel EventCode = 490
	EvPickupCarriableObjectFinished EventCode = 491
	EvDoSimpleActionStart EventCode = 492
	EvDoSimpleActionCancel EventCode = 493
	EvDoSimpleActionFinished EventCode = 494
	EvNotifyGuestAccountVerified EventCode = 495
	EvMightAndFavorReceivedEvent EventCode = 496
	EvWeeklyPvpChallengeRewardStateUpdate EventCode = 497
	EvNewUnlockedPvpSeasonChallengeRewards EventCode = 498
	EvStaticDungeonEntrancesDungeonEventStatusUpdates EventCode = 499
	EvStaticDungeonDungeonValueUpdate EventCode = 500
	EvStaticDungeonEntranceDungeonEventsAborted EventCode = 501
	EvInAppPurchaseConfirmedGooglePlay EventCode = 502
	EvFeatureSwitchInfo EventCode = 503
	EvPartyJoinRequestAborted EventCode = 504
	EvPartyInviteAborted EventCode = 505
	EvPartyStartHuntRequest EventCode = 506
	EvPartyStartHuntRequested EventCode = 507
	EvPartyStartHuntRequestAnswer EventCode = 508
	EvPartyPlayerLeaveScheduled EventCode = 509
	EvGuildInviteDeclined EventCode = 510
	EvCancelMultiSpellSlots EventCode = 511
	EvNewVisualEventObject EventCode = 512
	EvCastleClaimProgress EventCode = 513
	EvCastleClaimProgressLogo EventCode = 514
	EvTownPortalUpdateState EventCode = 515
	EvTownPortalFailed EventCode = 516
	EvConsumableVanityChargesAdded EventCode = 517
	EvFestivitiesUpdate EventCode = 518
	EvNewBannerObject EventCode = 519
	EvNewMistsImmediateReturnExit EventCode = 520
	EvMistsPlayerJoinedInfo EventCode = 521
	EvNewMistsStaticEntrance EventCode = 522
	EvNewMistsOpenWorldExit EventCode = 523
	EvNewTunnelExitTemp EventCode = 524
	EvNewMistsWispSpawn EventCode = 525
	EvMistsWispSpawnStateChange EventCode = 526
	EvNewMistsCityEntrance EventCode = 527
	EvNewMistsCityRoadsEntrance EventCode = 528
	EvMistsCityRoadsEntrancePartyStateUpdate EventCode = 529
	EvMistsCityRoadsEntranceClearStateForParty EventCode = 530
	EvMistsEntranceDataChanged EventCode = 531
	EvNewCagedObject EventCode = 532
	EvCagedObjectStateUpdated EventCode = 533
	EvEntrancePartyBindingCreated EventCode = 534
	EvEntrancePartyBindingCleared EventCode = 535
	EvEntrancePartyBindingInfos EventCode = 536
	EvNewMistsBorderExit EventCode = 537
	EvNewMistsDungeonExit EventCode = 538
	EvLocalQuestInfos EventCode = 539
	EvLocalQuestStarted EventCode = 540
	EvLocalQuestActive EventCode = 541
	EvLocalQuestInactive EventCode = 542
	EvLocalQuestProgressUpdate EventCode = 543
	EvNewUnrestrictedPvpZone EventCode = 544
	EvTemporaryFlaggingStatusUpdate EventCode = 545
	EvSpellTestPerformanceUpdate EventCode = 546
	EvTransformation EventCode = 547
	EvTransformationEnd EventCode = 548
	EvUpdateTrustlevel EventCode = 549
	EvRevealHiddenTimeStamps EventCode = 550
	EvModifyItemTraitFinished EventCode = 551
	EvRerollItemTraitValueFinished EventCode = 552
	EvHuntQuestProgressInfo EventCode = 553
	EvHuntStarted EventCode = 554
	EvHuntFinished EventCode = 555
	EvHuntAborted EventCode = 556
	EvHuntMissionStepStateUpdate EventCode = 557
	EvNewHuntTrack EventCode = 558
	EvHuntMissionUpdate EventCode = 559
	EvHuntQuestMissionProgressUpdate EventCode = 560
	EvHuntTrackUsed EventCode = 561
	EvHuntTrackUseableAgain EventCode = 562
	EvMinimapHuntTrackMarkers EventCode = 563
	EvNoTracksFound EventCode = 564
	EvHuntQuestAborted EventCode = 565
	EvInteractWithTrackStart EventCode = 566
	EvInteractWithTrackCancel EventCode = 567
	EvInteractWithTrackFinished EventCode = 568
	EvNewDynamicCompound EventCode = 569
	EvLegendaryItemDestroyed EventCode = 570
	EvAttunementInfo EventCode = 571
	EvTerritoryClaimRaidedRawEnergyCrystalResult EventCode = 572
	EvCarriedObjectExpiryWarning EventCode = 573
	EvCarriedObjectExpired EventCode = 574
	EvTerritoryRaidStart EventCode = 575
	EvTerritoryRaidCancel EventCode = 576
	EvTerritoryRaidFinished EventCode = 577
	EvTerritoryRaidResult EventCode = 578
	EvTerritoryMonolithActiveRaidStatus EventCode = 579
	EvTerritoryMonolithActiveRaidCancelled EventCode = 580
	EvMonolithEnergyStorageUpdate EventCode = 581
	EvMonolithNextScheduledOpenWorldAttackUpdate EventCode = 582
	EvMonolithProtectedBuildingsDamageReductionUpdate EventCode = 583
	EvNewBuildingBaseEvent EventCode = 584
	EvNewFortificationBuilding EventCode = 585
	EvNewCastleGateBuilding EventCode = 586
	EvBuildingDurabilityUpdate EventCode = 587
	EvMonolithFortificationPointsUpdate EventCode = 588
	EvFortificationBuildingUpgradeInfo EventCode = 589
	EvFortificationBuildingsDamageStateUpdate EventCode = 590
	EvSiegeNotificationEvent EventCode = 591
	EvUpdateEnemyWarBannerActive EventCode = 592
	EvTerritoryAnnouncePlayerEjection EventCode = 593
	EvCastleGateSwitchUseStarted EventCode = 594
	EvCastleGateSwitchUseFinished EventCode = 595
	EvFortificationBuildingWillDowngrade EventCode = 596
	EvBotCommand EventCode = 597
	EvJournalAchievementProgressUpdate EventCode = 598
	EvJournalClaimableRewardUpdate EventCode = 599
	EvKeySync EventCode = 600
	EvLocalQuestAreaGone EventCode = 601
	EvDynamicTemplate EventCode = 602
	EvDynamicTemplateForcedStateChange EventCode = 603
	EvNewOutlandsTeleportationPortal EventCode = 604
	EvNewOutlandsTeleportationReturnPortal EventCode = 605
	EvOutlandsTeleportationBindingCleared EventCode = 606
	EvOutlandsTeleportationReturnPortalUpdateEvent EventCode = 607
	EvPlayerUsedOutlandsTeleportationPortal EventCode = 608
	EvEncumberedRestricted EventCode = 609
	EvNewPiledObject EventCode = 610
	EvPiledObjectStateChanged EventCode = 611
	EvNewSmugglerCrateDeliveryStation EventCode = 612
	EvKillRewardedNoFame EventCode = 613
	EvPickupFromPiledObjectStart EventCode = 614
	EvPickupFromPiledObjectCancel EventCode = 615
	EvPickupFromPiledObjectReset EventCode = 616
	EvPickupFromPiledObjectFinished EventCode = 617
	EvArmoryActivityChange EventCode = 618
	EvNewKillTrophyFurnitureBuilding EventCode = 619
	EvHellDungeonsPlayerJoinedInfo EventCode = 620
	EvNewTileSwitchTrigger EventCode = 621
	EvNewMultiRewardObject EventCode = 622
	EvNewHellDungeonSoulShrineObject EventCode = 623
	EvHellDungeonSoulShrineStateUpdate EventCode = 624
	EvNewResurrectionShrine EventCode = 625
	EvUpdateResurrectionShrine EventCode = 626
	EvStandTimeFinished EventCode = 627
	EvEpicAchievementAndStatsUpdate EventCode = 628
	EvSpectateTargetAfterDeathUpdate EventCode = 629
	EvSpectateTargetAfterDeathEnded EventCode = 630
	EvNewHellDungeonUpwardExit EventCode = 631
	EvNewHellDungeonSoulExit EventCode = 632
	EvNewHellDungeonDownwardExit EventCode = 633
	EvNewHellDungeonChestExit EventCode = 634
	EvNewCorruptedStaticEntrance EventCode = 635
	EvNewHellDungeonStaticEntrance EventCode = 636
	EvUpdateHellDungeonStaticEntranceState EventCode = 637
	EvDebugTriggerHellDungeonShutdownStart EventCode = 638
	EvFullJournalQuestInfo EventCode = 639
	EvJournalQuestProgressInfo EventCode = 640
	EvNewHellDungeonRoomShrineObject EventCode = 641
	EvHellDungeonRoomShrineStateUpdate EventCode = 642
	EvSimpleBehaviourBuildingStateUpdate EventCode = 643
	EvSetTimeScaling EventCode = 644
	EvStopTimeScaling EventCode = 645
	EvKeyValidation EventCode = 646
	EvPlayerJoinMapMarkerTimerStates EventCode = 647
	EvNewMapMarkerTimer EventCode = 648
	EvRemoveMapMarkerTimer EventCode = 649
	EvNewFactionFortressObject EventCode = 650
	EvFactionFortressAnnouncePlayerEjection EventCode = 651
	EvRewardFactionWarfareSupply EventCode = 652
	EvFactionCaptureAreaProgressUpdate EventCode = 653
	EvFactionFortressClaimed EventCode = 654
	EvFactionFortressWeaponCachesSpawned EventCode = 655
	EvFactionFortressWeaponCacheClaimed EventCode = 656
	EvFactionFortressFightStateUpdate EventCode = 657
	EvFactionFortressCutoffFightStateUpdate EventCode = 658
	EvFactionFortressFightEnded EventCode = 659
	EvNewFactionWarfarePortal EventCode = 660
	EvFactionPortalTargetUpdate EventCode = 661
	EvFactionFortressFightStartedInRemoteClusterEvent EventCode = 662
	EvFactionFortressFightFinishedInRemoteClusterEvent EventCode = 663
	EvFactionDuchySupplyWarDefensiveVictoryEvent EventCode = 664
	EvFactionDuchyReconnectedFromCutoffEvent EventCode = 665
	EvFactionFortressCutoffFightCancelledByClusterOwnerChangeEvent EventCode = 666
	EvFactionDuchyEnteredCutoffStateEvent EventCode = 667
	EvLeaveProtectionStateUpdate EventCode = 668
	EvRedZoneEventStandings EventCode = 669
	EvNewFactionBattleStandardDeliveryStation EventCode = 670
	EvNewLoreSnippetObject EventCode = 671
	EvLoreSnippetObjectStateUpdate EventCode = 672
	EvLoreSnippedClaimed EventCode = 673
	EvLoreSnippetStatesChangedByCheat EventCode = 674
	EvNewTeleporterNode EventCode = 675
	EvTeleporterNodeStateChanged EventCode = 676
	EvTeleporterConnectionsFullStateUpdate EventCode = 677
	EvTeleporterConnectionStateChanged EventCode = 678
	EvRetrieveCarriableObjectStart EventCode = 679
	EvRetrieveCarriableObjectCancel EventCode = 680
	EvRetrieveCarriableObjectReset EventCode = 681
	EvRetrieveCarriableObjectFinished EventCode = 682
	EvLosingCarriableObjectStart EventCode = 683
	EvLosingCarriableObjectFinished EventCode = 684
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

