CREATE TABLE IF NOT EXISTS Merchants (
    MerchantId INT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS Offers (
    OfferId INT NOT NULL,
    OfferName VARCHAR(100) NOT NULL,
    Price INT NOT NULL,
    Quantity INT NOT NULL,
    MerchantId INT NOT NULL,
    CONSTRAINT compound PRIMARY KEY(OfferId, MerchantId),
    FOREIGN KEY (MerchantId) REFERENCES Merchants (MerchantId)
    ON UPDATE CASCADE
);