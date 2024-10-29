import { ReactElement } from 'react'
import {AppShell, Autocomplete, Avatar, Group, Loader, Menu, Stack, Text, UnstyledButton} from '@mantine/core'
import demoUser from '../../assets/demoUser.png'
import {IconCaretDownFilled, IconLogout2, IconSearch, IconSettings, IconUser} from '@tabler/icons-react'
import { useAppDispatch } from '../../state/store.ts'
import { signOut } from '../../state/authSlice.ts'
import { useGetCurrentUserQuery } from '../../state/api.ts'

function Topbar(): ReactElement {
  const dispatch = useAppDispatch()

  const { data: user, isLoading } = useGetCurrentUserQuery()

  function handleSignOut() {
    dispatch(signOut())
  }

  return (
    <AppShell.Header px={'md'} withBorder={false} top={'unset'}>
      <Group justify={'space-between'} align={'center'} h={'100%'}>
        <Autocomplete
          placeholder="Search"
          leftSection={<IconSearch size={16} stroke={1.5} />}
          data={[]}
          visibleFrom="xs"
          radius={'lg'}
          w={200}
        />

        {isLoading ? (
          <Loader />
        ) : (
          <Menu shadow={'lg'} width={200}>
            <Menu.Target>
              <UnstyledButton style={{ cursor: 'pointer' }} data-testid={'user-button'}>
                <Group gap={4}>
                  <Avatar src={demoUser} />
                  <IconCaretDownFilled size={12} color={'#323233'} />
                </Group>
              </UnstyledButton>
            </Menu.Target>

            <Menu.Dropdown>
              <Menu.Label>
                <Stack gap={0}>
                  <Text fw={400} c={'black'}>{user.name}</Text>
                  <Text fz={'xs'} fw={300}>{user.email}</Text>
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
