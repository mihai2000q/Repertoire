import { createTheme, rem } from '@mantine/core'
import { colors } from './colors'
import { components } from './components'
import classes from './active.module.css'

export const theme = createTheme({
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  colors: colors,
  primaryColor: 'primary',
  primaryShade: 4,
  white: '#fafafa',
  black: '#242424',

  defaultRadius: 'md',
  cursorType: 'pointer',

  activeClassName: classes.active,

  autoContrast: true,
  luminanceThreshold: 0.5,

  shadows: {
    xxl: 'rgba(0, 0, 0, 0.2) 0px 10px 36px 0px',
    xxl_hover: 'rgba(0, 0, 0, 0.4) 0px 10px 36px 0px'
  },

  // typography
  fontFamily: 'Poppins, sans-serif',
  fontSizes: {
    xxs: rem(9),
    xs: rem(10),
    sm: rem(11),
    md: rem(12),
    lg: rem(14),
    xl: rem(16)
  },
  lineHeights: {
    xxs: '1.15',
    xs: '1.3',
    sm: '1.4',
    md: '1.55',
    lg: '1.6',
    xl: '1.7'
  },
  breakpoints: {
    xs: '36em',
    sm: '48em',
    betweenSmMd: '54em',
    md: '62em',
    betweenMdLg: '68em',
    lg: '75em',
    betweenLgXl: '80em',
    xl: '88em',
    betweenXlXxl: '95em',
    xxl: '105em',
  },
  spacing: {
    xxs: rem(4)
  },

  headings: {
    fontWeight: '500',
    fontFamily: 'Poppins, sans-serif',
    sizes: {
      h1: {
        fontSize: rem(40),
        fontWeight: '900'
      },
      h2: {
        fontSize: rem(36),
        fontWeight: '800'
      },
      h3: {
        fontSize: rem(32),
        fontWeight: '700'
      },
      h4: {
        fontSize: rem(28),
        fontWeight: '600'
      },
      h5: { fontSize: rem(24) },
      h6: { fontSize: rem(20) }
    }
  },
  components: components
})
