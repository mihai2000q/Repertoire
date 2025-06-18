import { Divider, Grid, Skeleton, Stack } from '@mantine/core'

function PlaylistDrawerLoader() {
  return (
    <Stack gap={'xs'} data-testid={'playlist-drawer-loader'}>
      <Skeleton radius={0} w={'100%'} h={330} />

      <Stack px={'md'} pb={'md'} gap={6}>
        <Skeleton w={'max(20vw, 220px)'} h={'max(2.5vw, 25px)'} />
        <Skeleton w={50} h={15} mt={'xxs'} />
        <Divider my={'xs'} />
        <Stack gap={'md'}>
          {Array.from({ length: 5 }).map((_, i) => (
            <Grid key={i} align={'center'} gutter={'xs'} px={'sm'}>
              <Grid.Col span={1}>
                <Skeleton radius={'sm'} w={18} h={18} />
              </Grid.Col>
              <Grid.Col span={1.2}>
                <Skeleton radius={'8x'} w={28} h={28} />
              </Grid.Col>
              <Grid.Col span={9.6}>
                <Skeleton radius={'sm'} w={100} h={15} />
              </Grid.Col>
            </Grid>
          ))}
        </Stack>
      </Stack>
    </Stack>
  )
}

export default PlaylistDrawerLoader
