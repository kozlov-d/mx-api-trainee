import pandas as pd
import numpy as np
import string
import random
np.random.seed()
rows = 200000
length = 7
alph = np.array(list(string.ascii_lowercase))
data = pd.DataFrame({"OfferID" : list(range(1, 1 + rows)),
                     "OfferName"  : [(''.join(random.choices(string.ascii_uppercase, k=length))) for _ in range(rows)],
                     "Price"  : np.random.randint(0, 1000, size = rows),
                     "Quantity"  : np.random.randint(1, 1000, size = rows),
                     "Available"  : np.random.choice([True, False], size = rows)
                     })
data.to_excel("offers.xlsx", header=False, index=False)