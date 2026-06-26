import { Skeleton, Stack } from '@mantine/core'

function HomeArtistsLoader() {
  return (
    <>
      {Array.from(Array(20)).map((_, i) => (
        <Stack key={i} gap={'sm'} align={'center'} w={'max(9vw, 140px)'}>
          <Skeleton
            radius={'50%'}
            h={'max(calc(9vw - 25px), 125px)'}
            w={'max(calc(9vw - 25px), 125px)'}
            style={(theme) => ({ boxShadow: theme.shadows.md })}
          />
          <Stack gap={0} align={'center'}>
            <Skeleton w={100} h={15} mb={4} />
          </Stack>
        </Stack>
      ))}
    </>
  )
}

export default HomeArtistsLoader
