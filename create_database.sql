USE test_trademan;

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `PAIR`;
DROP TABLE IF EXISTS `WALLET`;
DROP TABLE IF EXISTS `POSITION`;
DROP TABLE IF EXISTS `ORDER`;
DROP TABLE IF EXISTS `LOG`;

DROP TRIGGER IF EXISTS positionModifyTrigger;

SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE `PAIR` (
	PairID INT NOT NULL AUTO_INCREMENT,
	Pair VARCHAR(15) NOT NULL UNIQUE,
	BaseCurrency VARCHAR(10),
	QuoteCurrency VARCHAR(10),
	PriceScale INT,
	TakerFee DECIMAL(20, 8),
	MakerFee DECIMAL(20, 8),
	MinLeverage DECIMAL(20, 8),
	MaxLeverage DECIMAL(20, 8),
	LeverageStep DECIMAL(20, 8),
	MinPrice DECIMAL(20, 8),
	MaxPrice DECIMAL(20, 8),
	TickSize DECIMAL(20, 8),
	MinOrderSize DECIMAL(20, 8),
	MaxOrderSize DECIMAL(20, 8),
	StepOrderSize DECIMAL(20, 8),

	INDEX(Pair),
	
	PRIMARY KEY (PairID)
);

CREATE TABLE `WALLET` (
	Symbol         VARCHAR(10) NOT NULL,
	Equity         DECIMAL(20, 8),
	Available      DECIMAL(20, 8),
	UsedMargin     DECIMAL(20, 8),
	OrderMargin    DECIMAL(20, 8),
	PositionMargin DECIMAL(20, 8),
	OCCClosingFee  DECIMAL(20, 8),
	OCCFundingFee  DECIMAL(20, 8),
	WalletBalance  DECIMAL(20, 8),
	DailyPnL       DECIMAL(20, 8),
	UnrealisedPnL  DECIMAL(20, 8),
	TotalPnL       DECIMAL(20, 8),
	EntryTime      DATETIME,
	
	PRIMARY KEY (Symbol, EntryTime)
);

CREATE TABLE POSITION (
	PositionID INT NOT NULL AUTO_INCREMENT,
	PairID INT NOT NULL,
	Status ENUM('Planned', 'Ordered', 'Filled', 'Stopped', 'Closed', 'Cancelled', 'Liquidated','Logged'),
	Side ENUM('Long', 'Short'),
	Risk DECIMAL(5,2),
	Notes TEXT,
	TradingViewPlan VARCHAR(100),
	RewardRiskRatio DECIMAL (6,2),
	Profit DECIMAL(21,12),
	EntryTime DATETIME DEFAULT CURRENT_TIMESTAMP,
	ModifyTime DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

	INDEX(PairID),
	INDEX(Status),
	INDEX(EntryTime),
	INDEX(ModifyTime),

	PRIMARY KEY (PositionID),
	FOREIGN KEY (PairID) REFERENCES `PAIR`(PairID)
);

CREATE TABLE `ORDER` (
	OrderID INT NOT NULL AUTO_INCREMENT,
	PositionID INT NOT NULL, 
	ExchangeOrderID VARCHAR(50),
	Status ENUM('Planned', 'Ordered', 'Position', 'Stopped', 'Closed', 'Cancelled', 'Logged'),
	OrderType ENUM('Hard Stoploss', 'Soft StopLoss', 'Entry', 'Take Profit'),
	`Size` DECIMAL(21,12),
	TriggerPrice DECIMAL(21,12),
	Price DECIMAL(21,12),
	EntryTime DATETIME DEFAULT CURRENT_TIMESTAMP,
	ModifyTime DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

	INDEX(EntryTime),
	INDEX(ModifyTime),

	PRIMARY KEY (OrderID),
	FOREIGN KEY (PositionID) REFERENCES `POSITION`(PositionID)
);

CREATE TABLE `LOG` (
	LogID INT NOT NULL AUTO_INCREMENT,
	PositionID INT NOT NULL,
	Source ENUM('Trigger', 'Software', 'User'),
	Text TEXT,
	EntryTime DATETIME DEFAULT CURRENT_TIMESTAMP,

	INDEX(EntryTime),
	INDEX(Source),

	PRIMARY KEY (LogID),
	FOREIGN KEY (PositionID) REFERENCES `POSITION`(PositionID)
);


DELIMITER $$

CREATE TRIGGER positionModifyTrigger
AFTER UPDATE ON `POSITION` FOR EACH ROW
BEGIN
	DECLARE logText Text;
	SET logText = "Updated fields: ";
	IF OLD.Status <> NEW.Status THEN
		SET logText = CONCAT(logText, 'Status: ', OLD.Status, '->', NEW.Status, ', ');
	END IF;
	IF OLD.Risk <> NEW.Risk THEN
		SET logText = CONCAT(logText, 'Risk: ', OLD.Risk, '->', NEW.Risk, ', ');
	END IF;
	IF OLD.Notes <> NEW.Notes THEN
		SET logText = CONCAT(logText, 'Notes: ', OLD.Notes, '->', NEW.Notes, ', ');
	END IF;
	IF OLD.TradingViewPlan <> NEW.TradingViewPlan THEN
		SET logText = CONCAT(logText, 'TradingViewPlan: ', OLD.TradingViewPlan, '->', NEW.TradingViewPlan, ', ');
	END IF;
	IF OLD.RewardRiskRatio <> NEW.RewardRiskRatio THEN
		SET logText = CONCAT(logText, 'Status: ', OLD.RewardRiskRatio, '->', NEW.RewardRiskRatio, ', ');
	END IF;
	IF OLD.Profit <> NEW.Profit THEN
		SET logText = CONCAT(logText, 'Profit: ', OLD.Profit, '->', NEW.Profit, ', ');
	END IF;
	INSERT INTO `LOG` (
		PositionID,
		Source,
		Text
	)
	VALUES (
		OLD.PositionID,
		'trigger',
		logText
	);
END$$

DELIMITER ;
