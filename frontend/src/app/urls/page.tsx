"use client";

import Link from "next/link";
import { useState } from "react";

const API_BASE = "http://localhost:8080";

export default function UrlListPage() {
  const [id, setId] = useState("");
  const [isFetching, setIsFetching] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [resolvedUrl, setResolvedUrl] = useState<string | null>(null);

  const handleFetch = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setResolvedUrl(null);
    setIsFetching(true);

    try {
      const response = await fetch(`${API_BASE}/${id}`, {
        method: "GET",
        redirect: "manual",
      });

      const location = response.headers.get("location");
      if (!location) {
        throw new Error("URLが取得できませんでした。");
      }

      setResolvedUrl(location);
    } catch (fetchError) {
      setError(
        fetchError instanceof Error
          ? fetchError.message
          : "URLの取得に失敗しました。"
      );
    } finally {
      setIsFetching(false);
    }
  };

  return (
    <div className="min-h-screen bg-zinc-50 px-6 py-12 text-zinc-900">
      <main className="mx-auto flex w-full max-w-2xl flex-col gap-8 rounded-2xl bg-white p-8 shadow-sm">
        <div className="flex items-start justify-between gap-6">
          <div>
            <h1 className="text-2xl font-semibold">URL一覧ページ</h1>
            <p className="mt-2 text-sm text-zinc-500">
              短縮IDから元のURLを取得します。
            </p>
          </div>
          <Link
            href="/"
            className="text-sm font-semibold text-blue-600 hover:text-blue-500"
          >
            Topへ戻る
          </Link>
        </div>

        <form className="flex flex-col gap-4" onSubmit={handleFetch}>
          <label className="text-sm font-medium text-zinc-700" htmlFor="id">
            ID
          </label>
          <input
            id="id"
            className="rounded-md border border-zinc-200 px-4 py-2 text-sm shadow-sm focus:border-zinc-400 focus:outline-none"
            placeholder="aBcD12eF"
            required
            value={id}
            onChange={(event) => setId(event.target.value)}
          />
          <button
            type="submit"
            className="rounded-md bg-zinc-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-zinc-800 disabled:cursor-not-allowed disabled:bg-zinc-400"
            disabled={!id || isFetching}
          >
            {isFetching ? "取得中..." : "取得"}
          </button>
        </form>

        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}

        {resolvedUrl && (
          <div className="rounded-md border border-zinc-200 bg-zinc-50 px-4 py-3 text-sm">
            <p className="font-medium text-zinc-700">取得結果</p>
            <p className="mt-2 break-all">
              <a
                href={resolvedUrl}
                className="text-blue-600 underline"
                target="_blank"
                rel="noreferrer"
              >
                {resolvedUrl}
              </a>
            </p>
          </div>
        )}
      </main>
    </div>
  );
}
