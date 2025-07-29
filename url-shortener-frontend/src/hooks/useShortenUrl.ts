/* eslint-disable @typescript-eslint/no-explicit-any */
import { useUrlStore } from "../store/UrlStore";
import { shortenAnon } from "../api/urls";
import type { UrlRequestAnon } from "../types/urls";

export function useShortenUrl() {
  const setLoading = useUrlStore((s) => s.setLoading);
  const setError = useUrlStore((s) => s.setError);
  const setAnonResult = useUrlStore((s) => s.setAnonResult);

  const shorten = async (payload: UrlRequestAnon) => {
    setLoading(true);
    setError(null);

    try {
      const result = await shortenAnon(payload);
      setAnonResult(result);
    } catch (err: any) {
      setError(err.message || "Something went wrong");
    } finally {
      setLoading(false);
    }
  };

  return { shorten };
}
