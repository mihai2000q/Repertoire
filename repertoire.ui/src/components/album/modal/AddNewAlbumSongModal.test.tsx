import { reduxRender } from '../../../test-utils.tsx'
import { setupServer } from 'msw/node'
import AddNewAlbumSongModal from './AddNewAlbumSongModal.tsx'

describe('Add New Album Song Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRender(<AddNewAlbumSongModal opened={true} onClose={() => {}} albumId={'album id'} />)
  })
})
