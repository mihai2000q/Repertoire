import {Divider, Grid, Group, Skeleton, Stack} from '@mantine/core'

function SongDrawerLoader() {
  return (
    <Stack gap={'xs'} data-testid={'song-drawer-loader'}>
      <Skeleton radius={0} w={'100%'} h={330} />

      <Stack px={'md'} pb={'xs'} gap={4}>
        <Skeleton w={220} h={25} />
        <Group gap={4}>
          <Skeleton radius={'50%'} w={35} h={35} />
          <Skeleton w={80} h={15} />
          <Skeleton w={110} h={15} />
          <Skeleton w={45} h={15} />
        </Group>

        <Stack gap={4} my={'xs'} px={'xs'}>
          <Skeleton w={'100%'} h={8} />
          <Skeleton w={'100%'} h={8} />
          <Skeleton w={50} h={8} />
        </Stack>

        <Divider />

        <Grid gutter={'md'} p={'xs'}>
          <Grid.Col span={4}>
            <Skeleton w={100} h={12} />
          </Grid.Col>
          <Grid.Col span={8}>
            <Skeleton w={175} h={15} />
          </Grid.Col>

          <Grid.Col span={4}>
            <Skeleton w={50} h={12} />
          </Grid.Col>
          <Grid.Col span={8}>
            <Skeleton w={70} h={15} />
          </Grid.Col>

          <Grid.Col span={4}>
            <Skeleton w={75} h={12} />
          </Grid.Col>
          <Grid.Col span={8}>
            <Skeleton w={90} h={15} />
          </Grid.Col>

          <Grid.Col span={4}>
            <Skeleton w={60} h={12} />
          </Grid.Col>
          <Grid.Col span={8}>
            <Skeleton w={35} h={15} />
          </Grid.Col>

          <Grid.Col span={4}>
            <Skeleton w={80} h={12} />
          </Grid.Col>
          <Grid.Col span={8}>
            <Skeleton w={200} h={15} />
          </Grid.Col>

          <Grid.Col span={4}>
            <Skeleton w={75} h={12} />
          </Grid.Col>
          <Grid.Col span={8}>
            <Skeleton w={200} h={15} />
          </Grid.Col>
        </Grid>

        <Divider my={4} />

        <Group gap={'xs'} style={{ alignSelf: 'flex-end' }}>
          <Skeleton radius={'sm'} w={23} h={23} />
          <Skeleton radius={'sm'} w={23} h={23} />
        </Group>
      </Stack>
    </Stack>
  )
}

export default SongDrawerLoader
