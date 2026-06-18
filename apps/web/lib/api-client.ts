const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export type SymbolDto = {
  id: string;
  code: string;
  market: string;
  baseAsset: string;
  quoteAsset: string;
  enabled: boolean;
};

export async function listSymbols(): Promise<SymbolDto[]> {
  const response = await fetch(`${API_BASE_URL}/api/symbols`, {
    next: { revalidate: 15 },
  });
  if (!response.ok) {
    throw new Error(`Failed to load symbols: ${response.status}`);
  }
  return response.json();
}

