#!/usr/bin/env python3
import requests
import json
from typing import Dict, Optional
import time
import random


class APITester:
    def __init__(self, base_url: str = "http://localhost:10800/api/v1"):
        self.base_url = base_url
        self.tokens = {
            'admin': None,
            'teacher': None,
            'student': None
        }
        self.test_users = {
            'admin': {'username': 'admin_test', 'password': 'admin123', 'role': 'admin', 'name': 'Admin User',
                      'email': 'admin@test.com'},
            'teacher': {'username': 'teacher_test', 'password': 'teacher123', 'role': 'teacher', 'name': 'Teacher User',
                        'email': 'teacher@test.com'},
            'student': {'username': 'student_test', 'password': 'student123', 'role': 'student', 'name': 'Student User',
                        'email': 'student@test.com'}
        }
        self.test_update_users = {
            'admin': {'username': 'admin_test', 'password': 'admin123', 'role': 'admin', 'name': 'Admin User New',
                      'email': 'admin_new@test.com'},
            'teacher': {'username': 'teacher_test', 'password': 'teacher123', 'role': 'teacher', 'name': 'Teacher User New',
                        'email': 'teacher_new@test.com'},
            'student': {'username': 'student_test', 'password': 'student123', 'role': 'student',
                        'name': 'Student User New',
                        'email': 'student_new@test.com'}
        }
        self.created_course_id = None
        self.enrollment_id = None

    def make_request(self, method: str, endpoint: str, data: Optional[Dict] = None,
                     token: Optional[str] = None) -> requests.Response:
        url = f"{self.base_url}{endpoint}"
        headers = {'Content-Type': 'application/json'}
        # Print the request details
        print(f"\nRequest: {method} {url}")

        if data:
            print(f"Request: {json.dumps(data, indent=2)}")

        if token:
            headers['Authorization'] = f'Bearer {token}'
        # print(f"Headers: {headers}")
        try:
            if method == 'GET':
                response = requests.get(url, headers=headers)
            elif method == 'POST':
                response = requests.post(url, json=data, headers=headers)
            elif method == 'PUT':
                response = requests.put(url, json=data, headers=headers)
            elif method == 'DELETE':
                response = requests.delete(url, headers=headers)
            else:
                raise ValueError(f"Unsupported HTTP method: {method}")

            print(f"\n{method} {endpoint}")
            print(f"Status Code: {response.status_code}")
            print(f"Response: {response.text[:200]}")
            return response
        except requests.exceptions.RequestException as e:
            print(f"Error making request: {e}")
            return None

    def register_user(self, user_type: str) -> bool:
        user_data = self.test_users[user_type]
        response = self.make_request('POST', '/auth/register', user_data)
        return response and response.status_code in [201, 200]

    def login_user(self, user_type: str) -> bool:
        user_data = {
            'username': self.test_users[user_type]['username'],
            'password': self.test_users[user_type]['password']
        }
        response = self.make_request('POST', '/auth/login', user_data)
        if response and response.status_code == 200:
            self.tokens[user_type] = response.json().get('token')
            return True
        return False

    def test_user_profile(self, user_type: str) -> bool:
        response = self.make_request('GET', '/users/profile', token=self.tokens[user_type])
        return response and response.status_code == 200

    def update_user_profile(self, user_type: str) -> bool:
        user_data = {
            'name': self.test_update_users[user_type]['name'],
            'email': self.test_update_users[user_type]['email']
        }
        response = self.make_request('PUT', '/users/profile', data=user_data, token=self.tokens[user_type])
        return response and response.status_code == 200

    def test_course_operations(self) -> bool:
        # Create course (as teacher)
        course_data = {
            'name': f'Test Course {int(time.time())}',
            'description': 'A test course for API testing',
            'credits': 3,
            'capacity': 30
        }
        response = self.make_request('POST', '/courses', course_data, self.tokens['teacher'])
        if not (response and response.status_code == 201):
            return False

        self.created_course_id = response.json().get('id')

        # List courses
        response = self.make_request('GET', '/courses', token=self.tokens['student'])
        if not (response and response.status_code == 200):
            return False

        # Get course details
        response = self.make_request('GET', f'/courses/{self.created_course_id}', token=self.tokens['student'])
        if not (response and response.status_code == 200):
            return False

        # Update course (as teacher)
        update_data = {
            'name': f'Updated Course {int(time.time())}',
            'description': 'Updated description'
        }
        response = self.make_request('PUT', f'/courses/{self.created_course_id}', update_data, self.tokens['teacher'])
        return response and response.status_code == 200

    def test_enrollment_operations(self) -> bool:
        if not self.created_course_id:
            print("No course available for enrollment testing")
            return False

        # Student enrolls in course
        response = self.make_request('POST', f'/courses/{self.created_course_id}/enroll', token=self.tokens['student'])
        if not (response and response.status_code in [200, 201]):
            return False

        # List enrollments (as student)
        response = self.make_request('GET', '/enrollments', token=self.tokens['student'])
        if not (response and response.status_code == 200):
            return False

        # Get student grades
        response = self.make_request('GET', '/enrollments/grades', token=self.tokens['student'])
        if not (response and response.status_code == 200):
            return False

        # Teacher updates grade
        if response.json().get('enrollments'):
            self.enrollment_id = response.json()['enrollments'][0]['id']
            grade_data = {
                'grade': round(random.uniform(60, 100), 2)
            }
            response = self.make_request('PUT', f'/enrollments/{self.enrollment_id}/grade',
                                         grade_data, self.tokens['teacher'])
            if not (response and response.status_code == 200):
                return False

        # Drop course (as student)
        response = self.make_request('POST', f'/courses/{self.created_course_id}/drop', token=self.tokens['student'])
        return response and response.status_code == 200

    def cleanup(self) -> bool:
        if self.created_course_id:
            # Delete course (as teacher)
            response = self.make_request('DELETE', f'/courses/{self.created_course_id}', token=self.tokens['teacher'])
            return response and response.status_code == 200
        return True

    def run_all_tests(self):
        test_results = {
            'user_registration': [],
            'user_login': [],
            'profile_access': [],
            'course_operations': None,
            'enrollment_operations': None,
            'cleanup': None
        }

        # Test user registration and login for each role
        for user_type in ['admin', 'teacher', 'student']:
            test_results['user_registration'].append({
                'user': user_type,
                'success': self.register_user(user_type)
            })
            test_results['user_login'].append({
                'user': user_type,
                'success': self.login_user(user_type)
            })
            test_results['profile_access'].append({
                'user': user_type,
                'success': self.test_user_profile(user_type)
            })

        # Test course operations
        test_results['course_operations'] = self.test_course_operations()

        # Test enrollment operations
        test_results['enrollment_operations'] = self.test_enrollment_operations()

        # Cleanup
        test_results['cleanup'] = self.cleanup()

        # Print results
        print("\nTest Results:")
        print(json.dumps(test_results, indent=2))

        # Calculate success rate
        total_tests = (len(test_results['user_registration']) +
                       len(test_results['user_login']) +
                       len(test_results['profile_access']) +
                       3)  # course_operations, enrollment_operations, cleanup

        successful_tests = (
                sum(1 for r in test_results['user_registration'] if r['success']) +
                sum(1 for r in test_results['user_login'] if r['success']) +
                sum(1 for r in test_results['profile_access'] if r['success']) +
                sum(1 for r in [test_results['course_operations'],
                                test_results['enrollment_operations'],
                                test_results['cleanup']] if r)
        )

        success_rate = (successful_tests / total_tests) * 100
        print(f"\nSuccess Rate: {success_rate:.2f}%")


if __name__ == "__main__":
    # Create API tester instance
    tester = APITester()

    user_type = 'student'
    # Run all tests
    tester.register_user(user_type)
    # time.sleep(3)
    tester.login_user(user_type)
    # tester.test_user_profile(user_type='student')
    tester.update_user_profile(user_type)
    tester.login_user(user_type)
