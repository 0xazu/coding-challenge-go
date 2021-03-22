import random
import uuid

from locust import HttpUser, task, between


# Load test to tet the service behaviour under concurrent load
class ChallengeLoadTest(HttpUser):
    wait_time = between(1, 2)
    users = []

    def on_start(self):
        # We generate 100 random user ids
        for i in range(0, 100):
            self.users.append(uuid.uuid4())

    @task(10)
    def post_transaction(self):
        # We generate a random small amount between 0 and 5 with two decimals
        amount = round(random.uniform(0.00, 5.00), 2)

        # We select a random user from the list
        user_id = self.users[random.randint(0, 100)]

        self.client.post('/transaction', headers={
            'Accept': 'application/json',
        }, json={
            'UserId': str(user_id),
            'Amount': amount
        }, timeout=10)

    @task(1)
    def get_batch(self):
        # We select a random user from the list
        user_id = self.users[random.randint(0, 100)]

        with self.client.get(f'/batch/{str(user_id)}', catch_response=True, timeout=10) as response:
            # We get a 404 because we didn't post yet the transaction for that user_id
            # We consider it a successful response
            if response.status_code == 404:
                response.success()

    @task(1)
    def get_batch_history(self):
        # We select a random user from the list
        user_id = self.users[random.randint(0, 100)]

        with self.client.get(f'/batch/history/{str(user_id)}', catch_response=True, timeout=10) as response:
            # We get a 404 because we didn't reach the threshold yet to move the batch to the history for that user_id
            # We consider it a successful response
            if response.status_code == 404:
                response.success()
