import { mantineRender } from '../../test-utils.tsx'
import PlaylistsLoader from './PlaylistsLoader.tsx'

describe('Playlists Loader', () => {
  it('should render', () => {
    mantineRender(<PlaylistsLoader />)
  })
})
