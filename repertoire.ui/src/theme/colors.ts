import { DefaultMantineColor, MantineColorsTuple } from '@mantine/core'

export const colors = {
  // cyan
  primary: [
    '#e3fafc',
    '#c5f6fa',
    '#99e9f2',
    '#66d9e8',
    '#3bc9db',
    '#22b8cf',
    '#15aabf',
    '#1098ad',
    '#0c8599',
    '#0b7285'
  ]
}

declare module '@mantine/core' {
  type ExtendedCustomColors = DefaultMantineColor | 'primary'

  // noinspection JSUnusedGlobalSymbols
  export interface MantineThemeColorsOverride {
    colors: Record<ExtendedCustomColors, MantineColorsTuple>
  }
}
