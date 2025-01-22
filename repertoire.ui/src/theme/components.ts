import {
  ActionIcon,
  ActionIconFactory,
  alpha,
  Button,
  Card,
  CardFactory,
  Combobox,
  LoadingOverlay,
  Menu,
  Modal,
  NavLink,
  NumberFormatter,
  Select,
  StylesApiProps,
  Text,
  Title,
  Tooltip,
  TooltipFloating
} from '@mantine/core'
import { DatePickerInput } from '@mantine/dates'

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
          },

          '&[data-disabled="true"]': {
            color: theme.colors.gray[3],
            backgroundColor: 'transparent'
          }
        },

        '&[data-variant="grey-subtle"]': {
          color: theme.colors.gray[2],
          backgroundColor: alpha(theme.colors.gray[6], 0.3),
          shadows: theme.shadows.xs,

          '&:hover': {
            color: theme.white,
            backgroundColor: alpha(theme.colors.gray[4], 0.3),
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
          color: theme.colors.primary[7],
          '&:hover': {
            boxShadow: theme.shadows.xxl_hover,
            color: theme.colors.primary[8],
            backgroundColor: alpha(theme.colors.primary[0], 0.2),
            transform: 'scale(1.1)'
          }
        }
      }
    })
  }),
  Combobox: Combobox.extend({
    defaultProps: {
      shadow: 'sm',
      transitionProps: {
        transition: 'scale-y',
        duration: 160
      }
    }
  }),
  DatePickerInput: DatePickerInput.extend({
    defaultProps: {
      valueFormat: 'D MMMM YYYY',
      dropdownType: 'popover',
      popoverProps: {
        shadow: 'sm',
        transitionProps: {
          transition: 'pop-top-left',
          duration: 160
        }
      }
    }
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
      },
      transitionProps: { transition: 'pop-top-right', duration: 150 }
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
      styles: (theme) => ({
        root: {
          borderRadius: '16px',
          transition: '0.25s',
          color: theme.colors.gray[6],
          '&:hover': {
            backgroundColor: 'inherit',
            color: theme.colors.gray[7],
            transform: 'scale(1.1)'
          },

          '&:where([data-active])': {
            backgroundColor: alpha(theme.colors.primary[4], 0.1),
            color: theme.colors.primary[4],

            '&:hover': {
              backgroundColor: alpha(theme.colors.primary[4], 0.15),
              color: theme.colors.primary[5],
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
      decimalSeparator: ',',
      decimalScale: 0
    }
  }),
  Select: Select.extend({
    defaultProps: {
      comboboxProps: {
        shadow: 'sm',
        transitionProps: {
          transition: 'scale-y',
          duration: 160
        }
      }
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
      bg: 'primary.9'
    }
  }),
  TooltipFloating: TooltipFloating.extend({
    defaultProps: {
      bg: 'primary.9'
    }
  })
}

declare module '@mantine/core' {
  // noinspection JSUnusedGlobalSymbols
  interface ActionIconProps {
    variant?: StylesApiProps<ActionIconFactory>['variant'] | 'grey' | 'grey-subtle'
  }

  // noinspection JSUnusedGlobalSymbols
  interface CardProps {
    variant?: StylesApiProps<CardFactory>['variant'] | 'panel' | 'add-new'
  }
}
