import {
  Anchor,
  Card,
  CardProps,
  Center,
  Grid,
  Group,
  ScrollArea,
  Stack,
  Text
} from '@mantine/core'

const mockGenres: { genre: string; color: string; isSmall: boolean }[] = [
  {
    genre: 'Rock',
    color: '#D9E7F2',
    isSmall: true
  },
  {
    genre: 'Techno',
    color: '#FEF9F3',
    isSmall: false
  },
  {
    genre: 'Indie',
    color: '#FFF0ED',
    isSmall: false
  },
  {
    genre: 'Hip Hop',
    color: '#E5EDEF',
    isSmall: true
  },
  {
    genre: 'Classical',
    color: '#FFF7FC',
    isSmall: true
  },
  {
    genre: 'Trap',
    color: '#F1EEFF',
    isSmall: false
  },
  {
    genre: 'Death Metal',
    color: '#fcf3ff',
    isSmall: false
  },
  {
    genre: 'Metal',
    color: '#f3fffb',
    isSmall: true
  },
  {
    genre: 'MetalCore',
    color: '#fff3f3',
    isSmall: true
  },
  {
    genre: 'Grunge',
    color: '#f5fff3',
    isSmall: false
  },
  {
    genre: 'Thrash Metal',
    color: '#fdfaec',
    isSmall: false
  },
  {
    genre: 'SoftCore',
    color: '#f2e1f6',
    isSmall: true
  }
]

function HomeGenres({ ...others }: CardProps) {
  return (
    <Card aria-label={'genres'} variant={'panel'} {...others} p={0}>
      <Stack gap={'xs'} h={'100%'}>
        <Group justify={'space-between'} px={'md'} pt={'md'}>
          <Text c={'gray.7'} fz={'lg'} fw={800}>
            Genres
          </Text>

          <Anchor
            fw={500}
            fz={'sm'}
            underline={'never'}
            c={'gray.5'}
            sx={{
              transition: '0.25s',
              '&:hover': { transform: 'scale(1.1)' },
              '&:active': { transform: 'scale(0.75)' }
            }}
          >
            See All
          </Anchor>
        </Group>

        <ScrollArea h={'100%'} scrollbars={'y'} scrollbarSize={7}>
          <Grid columns={10} px={'md'} pt={'sm'} gutter={'sm'}>
            {mockGenres.map((genre) => (
              <Grid.Col span={genre.isSmall ? 4 : 6} key={genre.genre}>
                <Card
                  key={genre.genre}
                  radius={'md'}
                  p={'lg'}
                  bg={genre.color}
                  sx={(theme) => ({
                    transition: '0.2s',
                    cursor: 'pointer',
                    boxShadow: theme.shadows.xs,
                    '&:hover': {
                      boxShadow: theme.shadows.sm,
                      transform: 'scale(1.05)'
                    },
                    '&:active': {
                      transform: 'scale(0.80)'
                    }
                  })}
                >
                  <Center>
                    <Text fw={600} c={'gray.7'} truncate={'end'}>
                      {genre.genre}
                    </Text>
                  </Center>
                </Card>
              </Grid.Col>
            ))}
          </Grid>
        </ScrollArea>
      </Stack>
    </Card>
  )
}

export default HomeGenres
