export interface ErrorDetail {
  field: string;
  message: string;
}

export interface APIErrorBody {
  code: string;
  message: string;
  details?: ErrorDetail[];
}

export interface APIErrorResponse {
  error: APIErrorBody;
}

export class APIError extends Error {
  constructor(
    public readonly code: string,
    message: string,
    public readonly details?: ErrorDetail[],
    public readonly status?: number,
  ) {
    super(message);
    this.name = "APIError";
  }
}

export interface ApiResponse<T> {
  data: T;
}
