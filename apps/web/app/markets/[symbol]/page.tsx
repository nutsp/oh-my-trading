import { MarketChartPage } from "@/features/markets/market-chart-page";

type MarketPageProps = {
  params: Promise<{ symbol: string }>;
  searchParams?: Promise<{ timeframe?: string }>;
};

export default async function MarketPage({ params, searchParams }: MarketPageProps) {
  const { symbol } = await params;
  const query = searchParams ? await searchParams : {};

  return (
    <MarketChartPage symbol={symbol.toUpperCase()} timeframe={query.timeframe ?? "1h"} />
  );
}
