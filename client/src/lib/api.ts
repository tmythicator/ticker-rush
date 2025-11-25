const API_URL = 'http://localhost:8080/api';

export interface Quote {
  symbol: string;
  price: number;
  timestamp: number;
}

export const fetchQuote = async (symbol: string): Promise<Quote> => {
  const res = await fetch(`${API_URL}/quote?symbol=${symbol}`);

  if (!res.ok) {
    throw new Error('Network response was not ok');
  }

  return res.json();
};