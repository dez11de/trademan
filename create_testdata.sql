USE test_trademan;

INSERT INTO `PLAN` (PairID, Status, Side, Risk, Notes, TradingViewPlan, RewardRiskRatio, Profit) 
VALUES((SELECT PairID FROM `PAIR` WHERE Pair = "BTCUSDT"), "statusPlanned", "sideLong", 0.5, "Created by create_testdata.sql.",
    "http://tradingview.com/someid/somenumber", 3.5, 0.0);

INSERT INTO `PLAN` (PairID, Status, Side, Risk, Notes, TradingViewPlan, RewardRiskRatio, Profit) 
VALUES(( SELECT PairID FROM `PAIR` WHERE Pair = "ETHUSDT"), "statusFilled", "sideLong", 0.5, "Created by create_testdata.sql.",
    "http://tradingview.com/someid/somedigits", 3.4, 0.0);

INSERT INTO `PLAN`  (PairID, Status, Side, Risk, Notes, TradingViewPlan, RewardRiskRatio, Profit) 
VALUES((SELECT PairID FROM `PAIR` WHERE Pair = "ADAUSDT"), "statusClosed", "sideShort", 1.0, "Created by create_testdata.sql.",
    "http://tradingview.com/someid/somenumber", 3.5, 12.33);

