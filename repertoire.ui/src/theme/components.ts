import {
  ActionIcon,
  ActionIconFactory,
  alpha,
  Button,
  Card,
  CardFactory,
  Chip,
  Combobox,
  Highlight,
  HoverCard,
  Indicator,
  LoadingOverlay,
  Menu,
  Modal,
  NavLink,
  NumberFormatter,
  NumberInput,
  RangeSlider,
  ScrollArea,
  ScrollAreaAutosize,
  Select,
  Stepper,
  StylesApiProps,
  Tabs,
  Text, Textarea,
  Title,
  Tooltip,
  TooltipFloating
} from '@mantine/core'
import { DatePickerInput } from '@mantine/dates'

export const components = {
  ActionIcon: ActionIcon.extend({
    styles: (theme) => ({
      root: {
        transition: '0.16s',
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
        },

        '&[data-variant="form"]': {
          color: theme.colors.gray[5],
          backgroundColor: theme.colors.gray[0],

          '&:hover': {
            color: theme.colors.gray[6],
            backgroundColor: theme.colors.gray[2]
          },

          '&[data-disabled="true"]': {
            color: theme.colors.gray[3],
            backgroundColor: 'transparent'
          },

          '&[aria-selected="true"]': {
            color: theme.colors.green[5],
            backgroundColor: alpha(theme.colors.green[1], 0.5),

            '&:hover': {
              color: theme.colors.green[6],
              backgroundColor: theme.colors.green[1]
            }
          },

          '&[aria-invalid="true"]': {
            border: `2px solid ${theme.colors.red[5]}`,
            color: theme.colors.red[7],
            backgroundColor: alpha(theme.colors.red[2], 0.75),

            '&:hover': {
              color: theme.colors.red[8],
              backgroundColor: theme.colors.red[2]
            }
          }
        }
      }
    })
  }),
  Button: Button.extend({
    defaultProps: {
      loaderProps: {
        type: 'dots'
      },
      styles: () => ({
        root: {
          transition: '0.18s'
        }
      })
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
          padding: 0,
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
  Chip: Chip.extend({
    styles: () => ({
      label: {
        transition: '0.15s',
        fontWeight: 500
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
      clearable: true,
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
  Highlight: Highlight.extend({
    defaultProps: {
      color: 'transparent',
      highlightStyles: { fontWeight: 700 }
    }
  }),
  HoverCard: HoverCard.extend({
    defaultProps: {
      shadow: 'md',
      withArrow: true
    }
  }),
  Indicator: Indicator.extend({
    defaultProps: {
      fw: 500
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
      content: {
        overflow: 'hidden'
      },
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
        body: { flex: 'unset' },
        section: {
          [`@media(min-width: ${theme.breakpoints.xxl})`]: {
            '& svg': {
              width: '26px',
              height: '26px'
            }
          }
        },
        label: {
          fontSize: theme.fontSizes.sm,
          fontWeight: 500,
          [`@media(max-width: ${theme.breakpoints.sm})`]: {
            fontSize: theme.fontSizes.lg,
            fontWeight: 600
          },
          [`@media(min-width: ${theme.breakpoints.xxl})`]: {
            fontSize: theme.fontSizes.md
          }
        },
        root: {
          [`@media(max-width: ${theme.breakpoints.sm})`]: {
            justifyContent: 'center',
            width: 'max(60vw, 200px)'
          },

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
  NumberInput: NumberInput.extend({
    defaultProps: {
      clampBehavior: 'strict'
    }
  }),
  RangeSlider: RangeSlider.extend({
    defaultProps: {
      minRange: 1
    }
  }),
  ScrollArea: ScrollArea.extend({
    defaultProps: {
      type: 'hover'
    }
  }),
  ScrollAreaAutosize: ScrollAreaAutosize.extend({
    defaultProps: {
      type: 'hover'
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
      },
      scrollAreaProps: {
        type: 'hover'
      }
    }
  }),
  Stepper: Stepper.extend({
    styles: (theme) => ({
      step: {
        borderRadius: theme.radius.md
      }
    })
  }),
  Tabs: Tabs.extend({
    defaultProps: {
      styles: (theme) => ({
        tab: {
          transition: 'border 200ms, background-color 200ms',
          '&:hover': {
            borderColor: theme.colors.gray[5],
            backgroundColor: alpha(theme.colors.gray[3], 0.6)
          },
          '&[data-active]': {
            '&:hover': {
              borderColor: theme.colors.primary[5]
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
  Textarea: Textarea.extend({
    styles: (theme) => ({
      input: {
        '&::-webkit-scrollbar': {
          width: 7
        },

        '&::-webkit-scrollbar-track-piece': {
          backgroundColor: 'transparent',
          '&:hover': {
            backgroundColor: theme.colors.gray[2]
          }
        },

        '&::-webkit-scrollbar-thumb': {
          borderRadius: theme.radius.md,
          backgroundColor: alpha(theme.colors.gray[6], 0.9),
          cursor: 'default',

          '&:hover': {
            backgroundColor: alpha(theme.colors.gray[7], 0.75)
          }
        }
      }
    })
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
    variant?: StylesApiProps<ActionIconFactory>['variant'] | 'grey' | 'grey-subtle' | 'form'
  }

  // noinspection JSUnusedGlobalSymbols
  interface CardProps {
    variant?: StylesApiProps<CardFactory>['variant'] | 'panel' | 'add-new'
  }
}
