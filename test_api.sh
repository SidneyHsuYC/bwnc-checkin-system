#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="http://localhost:8090"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🧪 Testing API Endpoints${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Test 1: Create User 1
echo -e "${YELLOW}📝 Test 1: Creating user - John Doe${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/user" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1-555-0101",
    "email": "john.doe@example.com"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 201 ]; then
  echo -e "${GREEN}✅ Success (HTTP $http_code)${NC}"
  echo "$body" | jq '.'
else
  echo -e "${RED}❌ Failed (HTTP $http_code)${NC}"
  echo "$body"
fi
echo ""

# Test 2: Create User 2
echo -e "${YELLOW}📝 Test 2: Creating user - Jane Smith${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/user" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "phone": "+1-555-0102",
    "email": "jane.smith@example.com"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 201 ]; then
  echo -e "${GREEN}✅ Success (HTTP $http_code)${NC}"
  echo "$body" | jq '.'
else
  echo -e "${RED}❌ Failed (HTTP $http_code)${NC}"
  echo "$body"
fi
echo ""

# Test 3: Create User 3
echo -e "${YELLOW}📝 Test 3: Creating user - Bob Johnson${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/user" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Bob",
    "last_name": "Johnson",
    "phone": "+1-555-0103",
    "email": "bob.johnson@example.com"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 201 ]; then
  echo -e "${GREEN}✅ Success (HTTP $http_code)${NC}"
  echo "$body" | jq '.'
else
  echo -e "${RED}❌ Failed (HTTP $http_code)${NC}"
  echo "$body"
fi
echo ""

# Test 4: Create User 4
echo -e "${YELLOW}📝 Test 4: Creating user - Alice Williams${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/user" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Alice",
    "last_name": "Williams",
    "phone": "+1-555-0104",
    "email": "alice.williams@example.com"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 201 ]; then
  echo -e "${GREEN}✅ Success (HTTP $http_code)${NC}"
  echo "$body" | jq '.'
else
  echo -e "${RED}❌ Failed (HTTP $http_code)${NC}"
  echo "$body"
fi
echo ""

# Test 5: Get All Users
echo -e "${YELLOW}📋 Test 5: Fetching all users${NC}"
response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/users")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
  echo -e "${GREEN}✅ Success (HTTP $http_code)${NC}"
  echo "$body" | jq '.'
else
  echo -e "${RED}❌ Failed (HTTP $http_code)${NC}"
  echo "$body"
fi
echo ""

# Test 6: Get User by ID
echo -e "${YELLOW}🔍 Test 6: Fetching user by ID (ID=1)${NC}"
response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/user/1")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
  echo -e "${GREEN}✅ Success (HTTP $http_code)${NC}"
  echo "$body" | jq '.'
else
  echo -e "${RED}❌ Failed (HTTP $http_code)${NC}"
  echo "$body"
fi
echo ""

# Test 7: Get User by ID - Not Found
echo -e "${YELLOW}🔍 Test 7: Fetching non-existent user (ID=9999)${NC}"
response=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/user/9999")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 404 ]; then
  echo -e "${GREEN}✅ Correctly returned 404 (HTTP $http_code)${NC}"
  echo "$body"
else
  echo -e "${RED}❌ Unexpected response (HTTP $http_code)${NC}"
  echo "$body"
fi
echo ""

# Test 8: Test validation - Missing fields
echo -e "${YELLOW}⚠️  Test 8: Testing validation (missing email)${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/user" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Invalid",
    "last_name": "User",
    "phone": "+1-555-9999"
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 400 ]; then
  echo -e "${GREEN}✅ Validation working correctly (HTTP $http_code)${NC}"
  echo "$body"
else
  echo -e "${RED}❌ Unexpected response (HTTP $http_code)${NC}"
  echo "$body"
fi
echo ""

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✨ Testing complete!${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
