CREATE TABLE IF NOT EXISTS Merchants (
    MerchantId INT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS Offers (
    Id SERIAL PRIMARY KEY,
    OfferId INT NOT NULL,
    OfferName VARCHAR(100) NOT NULL,
    Price INT NOT NULL,
    Quantity INT NOT NULL,
    CHECK((Price > 0) AND (Quantity > 0) AND (OfferName !='')),
    MerchantId INT NOT NULL,
    FOREIGN KEY (MerchantId) REFERENCES Merchants (MerchantId)
);