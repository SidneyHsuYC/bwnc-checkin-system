#!/bin/bash

# API Testing Script for Student Check-in System
# Run this script to test all the new features

API_URL="http://localhost:8090"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Unique suffix per run so the script is re-runnable: class_name and email are
# UNIQUE, so without this a second run would hit 409 conflicts on every create.
RUN_ID="$(date +%H%M%S)-$RANDOM"
LEADER_EMAIL="test.leader+${RUN_ID}@example.com"
NOCLASS_EMAIL="noclass+${RUN_ID}@example.com"
WITHCLASS_EMAIL="withclass+${RUN_ID}@example.com"

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║     Student Check-in System - API Test Suite                ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

# Check if server is running
echo -e "${BLUE}Checking server health...${NC}"
HEALTH=$(curl -s ${API_URL}/health)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Server is running${NC}"
    echo "$HEALTH" | jq '.'
else
    echo -e "${RED}✗ Server is not running. Please start the server first.${NC}"
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${YELLOW}TEST 1: Class Management${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo ""
echo -e "${BLUE}1.1 Creating a class (Monday 7pm-9pm)...${NC}"
CLASS1=$(curl -s -X POST ${API_URL}/api/classes \
  -H "Content-Type: application/json" \
  -d "{
    \"class_name\": \"Monday Test Class ${RUN_ID}\",
    \"start_date\": \"2026-02-17T00:00:00Z\",
    \"day_of_week\": \"Mon\",
    \"start_time\": \"19:00\",
    \"end_time\": \"21:00\"
  }")
CLASS1_ID=$(echo $CLASS1 | jq -r '.id')
if [ "$CLASS1_ID" != "null" ]; then
    echo -e "${GREEN}✓ Class created with ID: $CLASS1_ID${NC}"
    echo "$CLASS1" | jq '.'
else
    echo -e "${RED}✗ Failed to create class${NC}"
fi

echo ""
echo -e "${BLUE}1.2 Creating a student (for class leader test)...${NC}"
STUDENT1=$(curl -s -X POST ${API_URL}/api/students \
  -H "Content-Type: application/json" \
  -d "{
    \"first_name\": \"Test\",
    \"last_name\": \"Leader\",
    \"email\": \"${LEADER_EMAIL}\"
  }")
STUDENT1_ID=$(echo $STUDENT1 | jq -r '.id')
if [ "$STUDENT1_ID" != "null" ]; then
    echo -e "${GREEN}✓ Student created with ID: $STUDENT1_ID${NC}"
    echo "$STUDENT1" | jq '.'
else
    echo -e "${RED}✗ Failed to create student${NC}"
fi

echo ""
echo -e "${BLUE}1.3 Creating a class with leader (Tuesday 6pm-8pm)...${NC}"
CLASS2=$(curl -s -X POST ${API_URL}/api/classes \
  -H "Content-Type: application/json" \
  -d "{
    \"class_name\": \"Tuesday Test Class ${RUN_ID}\",
    \"start_date\": \"2026-02-18T00:00:00Z\",
    \"day_of_week\": \"Tue\",
    \"start_time\": \"18:00\",
    \"end_time\": \"20:00\",
    \"student_id\": $STUDENT1_ID
  }")
CLASS2_ID=$(echo $CLASS2 | jq -r '.id')
if [ "$CLASS2_ID" != "null" ]; then
    echo -e "${GREEN}✓ Class with leader created with ID: $CLASS2_ID${NC}"
    echo "$CLASS2" | jq '.'
else
    echo -e "${RED}✗ Failed to create class with leader${NC}"
fi

echo ""
echo -e "${BLUE}1.4 Listing all classes...${NC}"
CLASSES=$(curl -s ${API_URL}/api/classes)
CLASS_COUNT=$(echo $CLASSES | jq '. | length')
echo -e "${GREEN}✓ Found $CLASS_COUNT classes${NC}"
echo "$CLASSES" | jq '.'

echo ""
echo -e "${BLUE}1.5 Searching for Monday classes...${NC}"
SEARCH_RESULT=$(curl -s "${API_URL}/api/classes/search?q=mon")
SEARCH_COUNT=$(echo $SEARCH_RESULT | jq '. | length')
echo -e "${GREEN}✓ Found $SEARCH_COUNT Monday classes${NC}"
echo "$SEARCH_RESULT" | jq '.'

echo ""
echo -e "${BLUE}1.6 Getting class with leader details...${NC}"
CLASS_DETAIL=$(curl -s ${API_URL}/api/classes/$CLASS2_ID)
LEADER_NAME=$(echo $CLASS_DETAIL | jq -r '.leader.first_name + " " + .leader.last_name')
echo -e "${GREEN}✓ Class leader: $LEADER_NAME${NC}"
echo "$CLASS_DETAIL" | jq '.'

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${YELLOW}TEST 2: Student Registration (Optional Class Info)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo ""
echo -e "${BLUE}2.1 Creating student WITHOUT class info...${NC}"
STUDENT2=$(curl -s -X POST ${API_URL}/api/students \
  -H "Content-Type: application/json" \
  -d "{
    \"first_name\": \"NoClass\",
    \"last_name\": \"Student\",
    \"email\": \"${NOCLASS_EMAIL}\"
  }")
STUDENT2_ID=$(echo $STUDENT2 | jq -r '.id')
if [ "$STUDENT2_ID" != "null" ]; then
    echo -e "${GREEN}✓ Student created WITHOUT class info - ID: $STUDENT2_ID${NC}"
    echo "$STUDENT2" | jq '.'
else
    echo -e "${RED}✗ Failed to create student without class info${NC}"
fi

echo ""
echo -e "${BLUE}2.2 Creating student WITH class info and class_id...${NC}"
STUDENT3=$(curl -s -X POST ${API_URL}/api/students \
  -H "Content-Type: application/json" \
  -d "{
    \"first_name\": \"WithClass\",
    \"last_name\": \"Student\",
    \"class_info\": \"Mon 19:00-21:00\",
    \"class_id\": $CLASS1_ID,
    \"email\": \"${WITHCLASS_EMAIL}\"
  }")
STUDENT3_ID=$(echo $STUDENT3 | jq -r '.id')
if [ "$STUDENT3_ID" != "null" ]; then
    echo -e "${GREEN}✓ Student created WITH class info - ID: $STUDENT3_ID${NC}"
    echo "$STUDENT3" | jq '.'
else
    echo -e "${RED}✗ Failed to create student with class info${NC}"
fi

echo ""
echo -e "${BLUE}2.3 Searching for students...${NC}"
STUDENT_SEARCH=$(curl -s "${API_URL}/api/students/search?q=student")
STUDENT_SEARCH_COUNT=$(echo $STUDENT_SEARCH | jq '. | length')
echo -e "${GREEN}✓ Found $STUDENT_SEARCH_COUNT students${NC}"
echo "$STUDENT_SEARCH" | jq '.'

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${YELLOW}TEST 3: Check-in Flow${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo ""
echo -e "${BLUE}3.1 Creating an event...${NC}"
EVENT=$(curl -s -X POST ${API_URL}/api/events \
  -H "Content-Type: application/json" \
  -d "{
    \"event_name\": \"Test Workshop ${RUN_ID}\",
    \"event_time\": \"2026-02-20T19:00:00Z\",
    \"event_type\": \"Workshop\"
  }")
EVENT_ID=$(echo $EVENT | jq -r '.id')
if [ "$EVENT_ID" != "null" ]; then
    echo -e "${GREEN}✓ Event created with ID: $EVENT_ID${NC}"
    echo "$EVENT" | jq '.'
else
    echo -e "${RED}✗ Failed to create event${NC}"
fi

echo ""
echo -e "${BLUE}3.2 Checking in student to event...${NC}"
CHECKIN=$(curl -s -X POST ${API_URL}/api/checkins \
  -H "Content-Type: application/json" \
  -d "{
    \"student_id\": $STUDENT2_ID,
    \"event_id\": $EVENT_ID
  }")
CHECKIN_ID=$(echo $CHECKIN | jq -r '.id')
if [ "$CHECKIN_ID" != "null" ]; then
    echo -e "${GREEN}✓ Student checked in - Checkin ID: $CHECKIN_ID${NC}"
    echo "$CHECKIN" | jq '.'
else
    echo -e "${RED}✗ Failed to check in student${NC}"
fi

echo ""
echo -e "${BLUE}3.3 Listing check-ins for event...${NC}"
CHECKINS=$(curl -s "${API_URL}/api/checkins?event_id=$EVENT_ID")
CHECKIN_COUNT=$(echo $CHECKINS | jq '. | length')
echo -e "${GREEN}✓ Found $CHECKIN_COUNT check-ins for this event${NC}"
echo "$CHECKINS" | jq '.'

echo ""
echo -e "${BLUE}3.4 Testing duplicate check-in prevention...${NC}"
DUPLICATE=$(curl -s -X POST ${API_URL}/api/checkins \
  -H "Content-Type: application/json" \
  -d "{
    \"student_id\": $STUDENT2_ID,
    \"event_id\": $EVENT_ID
  }")
ERROR=$(echo $DUPLICATE | jq -r '.error')
if [ "$ERROR" != "null" ]; then
    echo -e "${GREEN}✓ Duplicate check-in prevented: $ERROR${NC}"
else
    echo -e "${RED}✗ Duplicate check-in was not prevented${NC}"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${YELLOW}TEST 4: Email Validation${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo ""
echo -e "${BLUE}4.1 Testing duplicate email prevention...${NC}"
DUPLICATE_EMAIL=$(curl -s -X POST ${API_URL}/api/students \
  -H "Content-Type: application/json" \
  -d "{
    \"first_name\": \"Duplicate\",
    \"last_name\": \"Email\",
    \"email\": \"${NOCLASS_EMAIL}\"
  }")
ERROR=$(echo $DUPLICATE_EMAIL | jq -r '.error')
if [[ "$ERROR" == *"already"* ]]; then
    echo -e "${GREEN}✓ Duplicate email prevented: $ERROR${NC}"
else
    echo -e "${RED}✗ Duplicate email was not prevented${NC}"
fi

echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                    TEST SUITE COMPLETE                       ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo -e "${GREEN}All tests completed!${NC}"
echo ""
echo "Summary:"
echo "  • Classes created: 2"
echo "  • Students created: 3"
echo "  • Events created: 1"
echo "  • Check-ins performed: 1"
echo ""
echo "You can now test the UI at: ${API_URL}"
echo ""
