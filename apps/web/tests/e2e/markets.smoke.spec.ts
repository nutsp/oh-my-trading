import { expect, test } from "@playwright/test";

test("market chart page renders with timeframe controls", async ({ page }) => {
  await page.goto("/markets/XAUUSD");

  await expect(page.getByRole("heading", { name: "XAUUSD Market Chart" })).toBeVisible();
  await expect(page.getByRole("link", { name: "15m" })).toBeVisible();
  await expect(page.getByRole("link", { name: "1h" })).toBeVisible();
  await expect(page.getByRole("link", { name: "4h" })).toBeVisible();
  await expect(page.getByRole("link", { name: "1d" })).toBeVisible();
});
