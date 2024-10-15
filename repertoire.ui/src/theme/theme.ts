import { createTheme, rem } from '@mantine/core'
import { colors } from './colors'

export const theme = createTheme({
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  colors: colors,
  primaryColor: 'cyan',
  primaryShade: { light: 5, dark: 7 },

  defaultRadius: 'md',
  cursorType: 'pointer',

  autoContrast: true,
  luminanceThreshold: 0.4,

  fontFamily: 'Open Sans, sans-serif',
  fontSizes: {
    xs: rem(10),
    sm: rem(11),
    md: rem(12),
    lg: rem(14),
    xl: rem(16)
  },
  lineHeights: {
    xs: '1.4',
    sm: '1.45',
    md: '1.55',
    lg: '1.6',
    xl: '1.65'
  },

  headings: {
    fontWeight: '500',
    fontFamily: 'Roboto, sans-serif',
    sizes: {
      h1: {
        fontSize: rem(40),
        fontWeight: '800'
      },
      h2: {
        fontSize: rem(36),
        fontWeight: '700'
      },
      h3: {
        fontSize: rem(32),
        fontWeight: '600'
      },
      h4: {
        fontSize: rem(28),
        fontWeight: '600'
      },
      h5: { fontSize: rem(24) },
      h6: { fontSize: rem(20) }
    }
  }
})
