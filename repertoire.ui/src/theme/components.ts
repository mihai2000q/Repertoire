import {
  ActionIcon,
  ActionIconFactory,
  alpha,
  Button,
  Card,
  CardFactory,
  LoadingOverlay,
  Menu,
  Modal,
  NavLink,
  NumberFormatter,
  StylesApiProps,
  Text,
  Title,
  Tooltip,
  TooltipFloating
} from '@mantine/core'

export const components = {
  ActionIcon: ActionIcon.extend({
    defaultProps: {
      style: { transition: '0.16s' }
    },
    styles: (theme) => ({
      root: {
        '&[data-variant="grey"]': {
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
  Button: Button.extend({
    defaultProps: {
      style: {
        transition: '0.18s'
      },
      loaderProps: {
        type: 'dots'
      }
    }
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
  LoadingOverlay: LoadingOverlay.extend({
    defaultProps: {
      overlayProps: { radius: 'md', blur: 2 },
      zIndex: 1000
    }
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
  Modal: Modal.extend({
    defaultProps: {
      closeButtonProps: {
        iconSize: 20
      }
    },
    styles: (theme) => ({
      title: {
        fontSize: theme.fontSizes.lg,
        fontWeight: 600,
        color: theme.colors.gray[7]
      }
    })
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
  NumberFormatter: NumberFormatter.extend({
    defaultProps: {
      thousandSeparator: ' ',
      decimalSeparator: ','
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
  }),
  TooltipFloating: TooltipFloating.extend({
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
