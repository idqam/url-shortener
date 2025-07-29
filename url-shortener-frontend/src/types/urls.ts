export type UrlStructure = {
  id?: string;
  original_url: string;
  short_url: string;
  short_code: string;
  user_id?: string | null;
  click_count?: number;
  created_at?: string;
};

export type UrlRequestAnon = {
  url: string;
  is_public: boolean;
  user_id: null;
  code_length: number;
};

export type UrlRequestAuthed = {
  url: string;
  is_public: boolean;
  user_id: string;
  code_length: number;
};

export type ShortenResponse = {
  id: string;
  short_code: string;
  original_url: string;
  short_url: string;
};

export type URLResponse = {
  id: string;
  user_id?: string | null;
  original_url: string;
  short_url: string;
  short_code: string;
  is_public: boolean;
  click_count: number;
  created_at: string;
};

export type GetUrlsResponse = {
  urls: URLResponse[];
};

export type UrlStructureWithAnalytics = UrlStructure & {
  click_count: number;
  analytics?: {
    country: string;
    clicks: number;
  }[];
};
