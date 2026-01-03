"use client";

import Link from "next/link";
import { useState } from "react";

const API_BASE = "http://localhost:8080";

export default function Home() {
  const [url, setUrl] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<{ id: string; shortUrl: string } | null>(
    null
  );

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setResult(null);
    setIsSubmitting(true);

    try {
      const response = await fetch(`${API_BASE}/shorten`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url }),
      });

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      const data = (await response.json()) as {
        id: string;
        short_url: string;
      };

      setResult({ id: data.id, shortUrl: data.short_url });
    } catch (submitError) {
      setError(
        submitError instanceof Error
          ? submitError.message
          : "Failed to create short URL."
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-zinc-50 px-6 py-12 text-zinc-900">
      <main className="mx-auto flex w-full max-w-2xl flex-col gap-8 rounded-2xl bg-white p-8 shadow-sm">
        <div>
          <h1 className="text-2xl font-semibold">URL Shortener</h1>
          <p className="mt-2 text-sm text-zinc-500">
            ローカルのAPI（{API_BASE}）に対してURLの短縮を行います。
          </p>
        </div>

        <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
          <label className="text-sm font-medium text-zinc-700" htmlFor="url">
            URL
          </label>
          <input
            id="url"
            className="rounded-md border border-zinc-200 px-4 py-2 text-sm shadow-sm focus:border-zinc-400 focus:outline-none"
            placeholder="https://example.com"
            type="url"
            required
            value={url}
            onChange={(event) => setUrl(event.target.value)}
          />
          <button
            type="submit"
            className="rounded-md bg-zinc-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-zinc-800 disabled:cursor-not-allowed disabled:bg-zinc-400"
            disabled={!url || isSubmitting}
          >
            {isSubmitting ? "作成中..." : "作成"}
          </button>
        </form>

        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}

        {result && (
          <div className="rounded-md border border-zinc-200 bg-zinc-50 px-4 py-3 text-sm">
            <p className="font-medium text-zinc-700">作成結果</p>
            <p className="mt-2">
              ID: <span className="font-mono">{result.id}</span>
            </p>
            <p>
              Short URL:{" "}
              <a
                href={result.shortUrl}
                className="text-blue-600 underline"
                target="_blank"
                rel="noreferrer"
              >
                {result.shortUrl}
              </a>
            </p>
          </div>
        )}

        <div>
          <Link
            href="/urls"
            className="inline-flex items-center text-sm font-semibold text-blue-600 hover:text-blue-500"
          >
            URL一覧ページへ
          </Link>
        </div>
      </main>
    </div>
  );
}
