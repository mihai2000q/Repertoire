import {
  ActionIcon,
  ActionIconFactory,
  alpha,
  Card,
  CardFactory,
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
  Card: Card.extend({
    styles: (theme) => ({
      root: {
        '&[data-variant="panel"]': {
          boxShadow: theme.shadows.sm,
          transition: '0.3s',
          '&:hover': {
            boxShadow: theme.shadows.xl
          }
        },
        '&[data-variant="add-new"]': {
          cursor: 'pointer',
          transition: '0.3s',
          boxShadow: theme.shadows.xxl,
          color: theme.colors.cyan[7],
          '&:hover': {
            boxShadow: theme.shadows.xxl_hover,
            color: theme.colors.cyan[8],
            backgroundColor: alpha(theme.colors.cyan[0], 0.2),
            transform: 'scale(1.1)'
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

  // noinspection JSUnusedGlobalSymbols
  interface CardProps {
    variant?: StylesApiProps<CardFactory>['variant'] | 'panel' | 'add-new'
  }
}
