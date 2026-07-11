export interface InvalidParam {
  name: string;
  reason: string;
}

export class ApiError extends Error {
  status: number;
  type?: string;
  title?: string;
  detail?: string;
  instance?: string;
  invalidParams?: InvalidParam[];

  constructor(
    message: string,
    status: number,
    details?: {
      type?: string;
      title?: string;
      detail?: string;
      instance?: string;
      invalidParams?: InvalidParam[];
    },
  ) {
    super(message);
    this.status = status;
    if (details) {
      this.type = details.type;
      this.title = details.title;
      this.detail = details.detail;
      this.instance = details.instance;
      this.invalidParams = details.invalidParams;
    }
  }
}

export interface ProblemDetailsPayload {
  type?: string;
  title?: string;
  detail?: string;
  instance?: string;
  invalid_params?: InvalidParam[];
  error?: string;
}

export const parseProblemDetails = (data: unknown, status: number): ApiError => {
  const payload = data as ProblemDetailsPayload | null;
  const detail = payload?.detail || payload?.error || `Error: ${status}`;
  return new ApiError(detail, status, {
    type: payload?.type,
    title: payload?.title,
    detail: payload?.detail,
    instance: payload?.instance,
    invalidParams: payload?.invalid_params,
  });
};
