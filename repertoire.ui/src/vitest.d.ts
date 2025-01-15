// noinspection JSUnusedGlobalSymbols

import 'vitest'

interface CustomMatchers<R = unknown> {
  toBeFormDataImage: (expected: File) => R
}

declare module 'vitest' {
  // eslint-disable-next-line @typescript-eslint/no-empty-object-type
  interface Assertion<T = never> extends CustomMatchers<T> {}

  // eslint-disable-next-line @typescript-eslint/no-empty-object-type
  interface AsymmetricMatchersContaining extends CustomMatchers {}
}
