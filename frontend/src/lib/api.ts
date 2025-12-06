const API_URL = import.meta.env.VITE_API_URL;

export interface Quote {
  symbol: string;
  price: number;
  timestamp: number;
}

class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

export const fetchQuote = async (symbol: string): Promise<Quote> => {
  const res = await fetch(`${API_URL}/quote?symbol=${symbol}`);

  if (!res.ok) {
    let errorMessage = `Error fetching quote: ${res.status}`;
    try {
      const errorData = await res.json();
      if (errorData.error) {
        errorMessage = errorData.error;
      }
    } catch {
      // failed to parse error json, use default message
    }

    throw new ApiError(errorMessage, res.status);
  }

  return res.json();
};

export interface User {
  user_id: number;
  email: string;
  balance: number;
  portfolio: Record<string, number>;
}

export const getUser = async (userId: number): Promise<User> => {
  const res = await fetch(`${API_URL}/user/${userId}`);
  if (!res.ok) throw new ApiError(`Error fetching user: ${res.status}`, res.status);
  return res.json();
};

export const buyStock = async (userId: number, symbol: string, count: number): Promise<User> => {
  const res = await fetch(`${API_URL}/buy`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user_id: userId, symbol, count }),
  });
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new ApiError(errorData.error || `Error buying stock: ${res.status}`, res.status);
  }
  return res.json();
};

export const sellStock = async (userId: number, symbol: string, count: number): Promise<User> => {
  const res = await fetch(`${API_URL}/sell`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user_id: userId, symbol, count }),
  });
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new ApiError(errorData.error || `Error selling stock: ${res.status}`, res.status);
  }
  return res.json();
};