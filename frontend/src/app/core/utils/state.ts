export interface State<T> {
  isLoading: boolean;
  value: T | null;
  error: string | null;
}

export function stateLoading<T>(): State<T> {
  return {
    isLoading: true,
    value: null,
    error: null,
  };
}

function getErrorMessage(err: unknown): string {
  if (typeof err === "string") return err;
  if (err instanceof Error) return err.message;
  return "An unknown error occurred";
}

export function stateError<T>(err: unknown): State<T> {
  return {
    isLoading: false,
    value: null,
    error: getErrorMessage(err),
  };
}

export function stateSuccess<T>(value: T): State<T> {
  return {
    isLoading: false,
    value: value,
    error: null,
  };
}
