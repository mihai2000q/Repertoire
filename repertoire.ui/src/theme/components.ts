import {
  ActionIcon,
  ActionIconFactory,
  alpha,
  Menu,
  NavLink,
  StylesApiProps,
  Text,
  Title,
  Tooltip
} from '@mantine/core'

export const components = {
  ActionIcon: ActionIcon.extend({
    styles: (theme) => ({
      root: {
        '&[data-variant="grey"]': {
          transition: '0.15s',
          color: theme.colors.gray[5],
          backgroundColor: theme.colors.gray[0],

          '&:hover': {
            color: theme.colors.gray[6],
            backgroundColor: theme.colors.gray[2],
            shadows: theme.shadows.lg
          }
        }
      }
    })
  }),
  Menu: Menu.extend({
    defaultProps: {
      styles: {
        item: {
          transition: '0.25s'
        }
      }
    }
  }),
  NavLink: NavLink.extend({
    defaultProps: {
      py: 'md',
      px: 'lg',
      style: {
        borderRadius: '16px'
      },
      styles: (theme) => ({
        root: {
          transition: '0.25s',
          color: theme.colors.gray[6],
          '&:hover': {
            backgroundColor: 'inherit',
            color: theme.colors.gray[7],
            transform: 'scale(1.1)'
          },

          '&:where([data-active])': {
            backgroundColor: alpha(theme.colors.cyan[4], 0.1),
            color: theme.colors.cyan[4],

            '&:hover': {
              backgroundColor: alpha(theme.colors.cyan[4], 0.15),
              color: theme.colors.cyan[5],
              transform: 'scale(1)'
            }
          }
        }
      })
    }
  }),
  Text: Text.extend({
    defaultProps: {
      c: 'dark'
    }
  }),
  Title: Title.extend({
    defaultProps: {
      c: 'dark'
    }
  }),
  Tooltip: Tooltip.extend({
    defaultProps: {
      bg: 'cyan.9'
    }
  })
}

declare module '@mantine/core' {
  // noinspection JSUnusedGlobalSymbols
  interface ActionIconProps {
    variant?: StylesApiProps<ActionIconFactory>['variant'] | 'grey'
  }
}
