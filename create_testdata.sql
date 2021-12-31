USE test_trademan;

INSERT INTO `POSITION` (Symbol, Status, Side, `Size`)
VALUES 	("BTCUSDT", "Planned", "long", 0.002),
		("ETHUSDT", "Stopped", "long", 1.23),
		("ADAUSDT", "Filled", "short", 54.21);
