import {alpha, NavLink} from "@mantine/core";

export const components = {
  NavLink: NavLink.extend({
    defaultProps: {
      py: 'md',
      px: 'lg',
      style: {
        borderRadius: '16px',
      },
      styles: (theme) => ({
        root: {
          transition: '0.25s',
          color: theme.colors.gray[6],
          '&:hover': {
            backgroundColor: 'inherit',
            color: theme.colors.gray[7],
            transform: 'scale(110%)',
          },

          '&:where([data-active])': {
            backgroundColor: alpha(theme.colors.cyan[4], 0.1),
            color: theme.colors.cyan[4],

            '&:hover': {
              backgroundColor: alpha(theme.colors.cyan[4], 0.15),
              color: theme.colors.cyan[5],
              transform: 'scale(100%)'
            },
          }
        },
      })
    },
  }),
}
