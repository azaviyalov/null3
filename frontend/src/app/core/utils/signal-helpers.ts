import { DestroyRef, WritableSignal, inject, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import {
  Observable,
  catchError,
  map,
  of,
  skip,
  startWith,
  switchMap,
} from "rxjs";
import {
  State as State,
  stateError,
  stateLoading,
  stateSuccess,
} from "./state";

export function toWritableSignal<TTrigger, TResult = TTrigger>({
  trigger,
  project = (input: TTrigger) => of(input as unknown as TResult),
  initialValue,
  onError = () => initialValue,
  destroyRef,
}: {
  trigger: Observable<TTrigger>;
  project?: (input: TTrigger) => Observable<TResult>;
  initialValue: TResult;
  onError?: (err: unknown) => TResult;
  destroyRef?: DestroyRef;
}): WritableSignal<TResult> {
  destroyRef = destroyRef ?? inject(DestroyRef);

  const result = signal<TResult>(initialValue);

  trigger
    .pipe(
      switchMap((input) => project(input).pipe(startWith(initialValue))),
      catchError((err) => {
        console.error(err);
        return of(onError(err));
      }),
      skip(1),
      takeUntilDestroyed(destroyRef),
    )
    .subscribe((value) => result.set(value));

  return result;
}

export function toWritableStateSignal<TTrigger, TResult = TTrigger>({
  trigger,
  project = (input: TTrigger) => of(input as unknown as TResult),
  destroyRef,
}: {
  trigger: Observable<TTrigger>;
  project?: (input: TTrigger) => Observable<TResult>;
  destroyRef?: DestroyRef;
}): WritableSignal<State<TResult>> {
  return toWritableSignal<TTrigger, State<TResult>>({
    trigger,
    project: (input: TTrigger) =>
      project(input).pipe(map((value) => stateSuccess(value))),
    initialValue: stateLoading(),
    onError: (err: unknown) => stateError<TResult>(err),
    destroyRef,
  });
}
