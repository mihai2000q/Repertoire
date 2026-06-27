import { mantineRender } from '../../../../../../test-utils.tsx'
import ArtistSongsLoader from './ArtistSongsLoader.tsx'

describe('Artist Songs Loader', () => {
  it('should render', () => {
    mantineRender(<ArtistSongsLoader />)
  })
})
