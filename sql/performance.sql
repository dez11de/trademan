SELECT (RecentEquity - PreviousEquity) / PreviousEquity * 100 AS Performance
FROM ( SELECT ( SELECT Equity
                FROM WALLET p1
                WHERE p1.EntryTime = x.PreviousTimestamp AND Symbol = x.Symbol
            ) AS PreviousEquity,
            ( SELECT Equity
                FROM WALLET p1
                WHERE p1.EntryTime = x.RecentTimestamp AND Symbol = x.Symbol
            ) AS RecentEquity
        FROM ( SELECT Symbol, MIN(EntryTime) AS PreviousTimestamp, MAX(EntryTime) AS RecentTimestamp
                FROM WALLET
                WHERE Symbol = "USDT" AND EntryTime BETWEEN '2022-01-27 00:00:00' AND NOW()
            ) x) x2;
