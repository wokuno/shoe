#include <WiFi.h>
#include <HTTPClient.h>

// WiFi network name and password:
const char * networkName = "UofM-Guest";
const char * networkPswd = "";

// Internet domain to request from:
const char * hostDomain = "esp.trana.fit";
const int hostPort = 80;

const char* serverNameTrial = "http://esp.trana.fit/newTrial";
const char* serverNameData = "http://esp.trana.fit/newData";

const int BUTTON_PIN = 0;
const int LED_PIN = 13;

int ledState = 0;
int trialID = -1;
int period = 20;
unsigned long time_now = 0;
int count = 0;
String jsonData = "[";

WiFiClient client;
HTTPClient http;
  
void setup()
{
  // Initilize hardware:
  Serial.begin(115200);
  pinMode(BUTTON_PIN, INPUT_PULLUP);
  pinMode(LED_PIN, OUTPUT);

  // Connect to the WiFi network (see function below loop)
  connectToWiFi(networkName, networkPswd);

  digitalWrite(LED_PIN, LOW); // LED off
  Serial.print("Press button 0 to connect to ");
  Serial.println(hostDomain);

}

void loop()
{

  if (digitalRead(BUTTON_PIN) == LOW)
  { // Check if button has been pressed
    while (digitalRead(BUTTON_PIN) == LOW)
      ; // Wait for button to be released
    if (trialID == -1) {
      trialID = httpGETRequest(serverNameTrial).toInt(); // Connect to server
      ledState = 1;
      time_now = millis();
    } else {
      uint16_t value1 = analogRead(32);
      uint16_t value2 = analogRead(33);
      uint16_t value3 = analogRead(34);
      uint16_t value4 = analogRead(39);
      uint16_t value5 = analogRead(36);
      jsonData += "{\"trialid\":" + String(trialID) + ",\"time\":" + String(millis()) + ",\"a1\":" + String(value1) + ",\"a2\":" + String(value2) + ",\"a3\":" + String(value3) + ",\"a4\":" + String(value4) + ",\"a5\":" + String(value5) + "}]";
//      Serial.println(jsonData);
      httpPOSTRequest(serverNameData, jsonData);
      count = 0;
      jsonData = "[";
      trialID = -1;
      ledState = 0;
    }
    digitalWrite(LED_PIN, ledState); // Turn off LED
    Serial.println(trialID);
  } else {
    uint16_t value1 = analogRead(32);
    uint16_t value2 = analogRead(33);
    uint16_t value3 = analogRead(34);
    uint16_t value4 = analogRead(39);
    uint16_t value5 = analogRead(36);
    if ((trialID != -1) && (millis() > time_now + period)) {
      if (count == 49) {
        jsonData += "{\"trialid\":" + String(trialID) + ",\"time\":" + String(millis()) + ",\"a1\":" + String(value1) + ",\"a2\":" + String(value2) + ",\"a3\":" + String(value3) + ",\"a4\":" + String(value4) + ",\"a5\":" + String(value5) + "}]";
        httpPOSTRequest(serverNameData, jsonData);
        count = 0;
        jsonData = "[";
      } else {
        jsonData += "{\"trialid\":" + String(trialID) + ",\"time\":" + String(millis()) + ",\"a1\":" + String(value1) + ",\"a2\":" + String(value2) + ",\"a3\":" + String(value3) + ",\"a4\":" + String(value4) + ",\"a5\":" + String(value5) + "},";
//        Serial.println(jsonData);
        count++;
      }
      time_now = millis();
    }
  }

}

void connectToWiFi(const char * ssid, const char * pwd)
{
  printLine();
  Serial.println("Connecting to WiFi network: " + String(ssid));

  WiFi.begin(ssid, pwd);

  while (WiFi.status() != WL_CONNECTED)
  {
    // Blink LED while we're connecting:
    digitalWrite(LED_PIN, ledState);
    ledState = (ledState + 1) % 2; // Flip ledState
    delay(500);
    Serial.print(".");
  }

  Serial.println();
  Serial.println("WiFi connected!");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());

}

String httpGETRequest(const char* serverName) {

  // Your Domain name with URL path or IP address with path
  http.begin(client, serverName);

  // Send HTTP POST request
  int httpResponseCode = http.GET();

  String payload = "--";

  if (httpResponseCode > 0) {
    Serial.print("HTTP Response code: ");
    Serial.println(httpResponseCode);
    payload = http.getString();
  }
  else {
    Serial.print("Error code: ");
    Serial.println(httpResponseCode);
  }
  // Free resources
  http.end();

  return payload;
}

int httpPOSTRequest(const char* serverName, String Data) {

  // Your Domain name with URL path or IP address with path
  http.begin(client, serverName);

  // If you need an HTTP request with a content type: application/json, use the following:
  http.addHeader("Content-Type", "application/json");
  int httpResponseCode = http.POST(Data);

  // Free resources
  http.end();

  return httpResponseCode;
}

void printLine()
{
  Serial.println();
  for (int i = 0; i < 30; i++)
    Serial.print("-");
  Serial.println();
}
