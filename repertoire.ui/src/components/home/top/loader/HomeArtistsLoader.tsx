import { Box, Skeleton, Stack } from '@mantine/core'

function HomeArtistsLoader() {
  return (
    <Box data-testid={'home-artists-loader'}>
      {Array.from(Array(20)).map((_, i) => (
        <Stack key={i} gap={'xs'} align={'center'}>
          <Skeleton
            radius={'50%'}
            h={125}
            w={125}
            style={(theme) => ({ boxShadow: theme.shadows.md })}
          />
          <Stack gap={0} align={'center'}>
            <Skeleton w={100} h={15} mb={4} />
          </Stack>
        </Stack>
      ))}
    </Box>
  )
}

export default HomeArtistsLoader
