import { Divider, Grid, Group, Skeleton, Stack } from '@mantine/core'

function SongLoader() {
  return (
    <Stack px={'xl'} data-testid={'song-loader'}>
      <Group>
        <Skeleton radius={'lg'} w={'max(12vw, 150px)'} h={'max(12vw, 150px)'} />
        <Stack gap={'xxs'} pt={'10px'}>
          <Skeleton w={65} h={15} />
          <Skeleton w={'max(25vw, 200px)'} h={'max(4vw, 48px)'} my={'xs'} />
          <Group gap={'xxs'}>
            <Skeleton radius={'50%'} w={35} h={35} />
            <Skeleton w={80} h={15} />
            <Skeleton w={110} h={15} />
            <Skeleton w={45} h={15} />
          </Group>
        </Stack>
      </Group>

      <Divider />

      <Grid align="start" mb={'lg'}>
        <Grid.Col span={{ sm: 12, md: 4.5 }}>
          <Stack>
            <Skeleton w={'100%'} h={220} />
            <Skeleton w={'100%'} h={140} />
            <Skeleton w={'100%'} h={130} />
          </Stack>
        </Grid.Col>

        <Grid.Col span={{ sm: 12, md: 7.5 }}>
          <Stack>
            <Skeleton w={'100%'} h={120} />
            <Skeleton w={'100%'} h={385} />
          </Stack>
        </Grid.Col>
      </Grid>
    </Stack>
  )
}

export default SongLoader
