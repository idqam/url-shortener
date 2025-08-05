export interface ShortenURLRequest {
  user_id?: string | null;
  url: string;
  is_public: boolean;
  code_length: number;
}

export interface ShortenURLResponse {
  id: string;
  short_code: string;
  short_url: string;
  created_at: string;
  is_public: boolean;
  click_count: number;
}

export interface GetUserURLsResponse {
  urls: ShortenURLResponse[];
}

export interface GetURLByShortCodeResponse {
  original_url: string;
  click_count: number;
}

export interface ErrorResponse {
  error: string;
  code?: string;
  field?: string;
}
