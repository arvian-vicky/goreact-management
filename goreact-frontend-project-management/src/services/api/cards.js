import network from '@/utils/network';

const cards = {
  async create(payload) {
    return network.post('/api/v1/cards', payload);
  },
  async update(id, payload) {
    return network.put(`/api/v1/cards/${id}`, payload);
  },
  async detail(id) {
    return network.get(`/api/v1/cards/${id}`);
  },
};

export default cards;
