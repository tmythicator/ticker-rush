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
    } catch (e) {
    }

    throw new ApiError(errorMessage, res.status);
  }

  return res.json();
};