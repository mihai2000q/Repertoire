import {NavLink} from "@mantine/core";

export const components = {
  NavLink: NavLink.extend({
    defaultProps: {
      py: 'md',
      px: 'lg',
      style: {
        borderRadius: '16px',
      }
    },
  }),
}
