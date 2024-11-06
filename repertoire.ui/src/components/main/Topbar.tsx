import { ReactElement } from 'react'
import {
  ActionIcon,
  alpha,
  AppShell,
  Autocomplete,
  Avatar,
  Group,
  Loader,
  Menu,
  Space,
  Stack,
  Text,
  UnstyledButton
} from '@mantine/core'
import userPlaceholder from '../../assets/user-placeholder.jpg'
import {
  IconBellFilled,
  IconCaretDownFilled,
  IconLogout2,
  IconSearch,
  IconSettings,
  IconUser
} from '@tabler/icons-react'
import { useAppDispatch } from '../../state/store.ts'
import { signOut } from '../../state/authSlice.ts'
import { useGetCurrentUserQuery } from '../../state/api.ts'
import useAuth from '../../hooks/useAuth.ts'

function Topbar(): ReactElement {
  const dispatch = useAppDispatch()

  const { data: user, isLoading } = useGetCurrentUserQuery(undefined, {
    skip: !useAuth()
  })

  function handleSignOut() {
    dispatch(signOut())
  }

  if (!user) return <></>

  return (
    <AppShell.Header px={'md'} withBorder={false} top={'unset'}>
      <Group align={'center'} h={'100%'} gap={0}>
        <Autocomplete
          placeholder="Search"
          leftSection={<IconSearch size={16} stroke={1.5} />}
          data={[]}
          visibleFrom="xs"
          radius={'lg'}
          w={200}
          styles={(theme) => ({
            input: {
              transition: '0.3s',
              backgroundColor: alpha(theme.colors.gray[0], 0.1),
              borderWidth: 0,

              '&:hover': {
                boxShadow: theme.shadows.sm,
                backgroundColor: alpha(theme.colors.gray[0], 0.2)
              },

              '&:focus': {
                boxShadow: theme.shadows.sm,
                backgroundColor: alpha(theme.colors.gray[0], 0.2)
              }
            }
          })}
        />

        <Space flex={1} />

        <ActionIcon
          variant={'subtle'}
          size={'lg'}
          styles={(theme) => ({
            root: {
              borderRadius: '50%',
              transition: '0.175s',
              color: theme.colors.gray[6],

              '&:hover': {
                backgroundColor: theme.colors.cyan[0],
                color: theme.colors.cyan[6]
              }
            }
          })}
        >
          <IconBellFilled size={18} />
        </ActionIcon>

        {isLoading ? (
          <Loader />
        ) : (
          <Menu shadow={'lg'} width={200}>
            <Menu.Target>
              <UnstyledButton
                p={'4px'}
                data-testid={'user-button'}
                sx={(theme) => ({
                  borderRadius: '16px',
                  cursor: 'pointer',
                  transition: '0.175s',
                  color: theme.colors.gray[7],
                  '&:hover': {
                    color: theme.colors.gray[8],
                    backgroundColor: alpha(theme.colors.gray[1], 0.7)
                  }
                })}
              >
                <Group gap={4}>
                  <Avatar src={user.imageUrl ? user.imageUrl : userPlaceholder} />
                  <IconCaretDownFilled size={12} />
                </Group>
              </UnstyledButton>
            </Menu.Target>

            <Menu.Dropdown>
              <Menu.Label>
                <Stack gap={0}>
                  <Text fw={400} c={'black'}>
                    {user.name}
                  </Text>
                  <Text fz={'xs'} fw={300}>
                    {user.email}
                  </Text>
                </Stack>
              </Menu.Label>

              <Menu.Divider />

              <Menu.Item leftSection={<IconUser size={14} />}>Account</Menu.Item>
              <Menu.Item leftSection={<IconSettings size={14} />}>Settings</Menu.Item>

              <Menu.Item leftSection={<IconLogout2 size={14} />} onClick={handleSignOut}>
                Sign Out
              </Menu.Item>
            </Menu.Dropdown>
          </Menu>
        )}
      </Group>
    </AppShell.Header>
  )
}

export default Topbar
