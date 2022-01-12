USE test_trademan;

INSERT INTO `POSITION` (PairID, Status, Side, Risk, Notes, TradingViewPlan, RewardRiskRatio, Profit) 
VALUES((SELECT PairID FROM `PAIR` WHERE Pair = "BTCUSDT"), "Planned", "Long", 0.5, "Created by create_testdata.sql.",
    "http://tradingview.com/someid/somenumber", 3.5, 0.0);

INSERT INTO `POSITION` (PairID, Status, Side, Risk, Notes, TradingViewPlan, RewardRiskRatio, Profit) 
VALUES(( SELECT PairID FROM `PAIR` WHERE Pair = "ETHUSDT"), "Filled", "Long", 0.5, "Created by create_testdata.sql.",
    "http://tradingview.com/someid/somedigits", 3.4, 0.0);

INSERT INTO `POSITION`  (PairID, Status, Side, Risk, Notes, TradingViewPlan, RewardRiskRatio, Profit) 
VALUES((SELECT PairID FROM `PAIR` WHERE Pair = "ADAUSDT"), "Closed", "Short", 1.0, "Created by create_testdata.sql.",
    "http://tradingview.com/someid/somenumber", 3.5, 12.33);
