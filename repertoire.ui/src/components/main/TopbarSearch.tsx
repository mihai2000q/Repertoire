import { IconSearch } from '@tabler/icons-react'
import { alpha, Autocomplete, AutocompleteProps } from '@mantine/core'

function TopbarSearch({ ...others }: AutocompleteProps) {
  return (
    <Autocomplete
      role={'searchbox'}
      aria-label={'topbar-search'}
      placeholder={'Search'}
      leftSection={<IconSearch size={16} stroke={2} />}
      data={[]}
      fw={500}
      radius={'lg'}
      styles={(theme) => ({
        input: {
          transition: '0.3s',
          backgroundColor: alpha(theme.colors.gray[0], 0.1),
          borderWidth: 0,
          '&:focus, &:hover': {
            boxShadow: theme.shadows.sm,
            backgroundColor: alpha(theme.colors.gray[0], 0.2)
          }
        }
      })}
      {...others}
    />
  )
}

export default TopbarSearch
